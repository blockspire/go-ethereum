// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package observer

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	lru "github.com/hashicorp/golang-lru"
)

// Trie cache generation limit after which to evic trie nodes from memory.
var MaxTrieCacheGen = uint16(120)

const (
	// maxPastTries defines the number of past tries to keep. This value is
	// chosen such that reasonable chain reorg depths will hit an existing trie.
	maxPastTries = 12

	// codeSizeCacheSize defines the number of codehash->size associations to keep.
	codeSizeCacheSize = 100000
)

// -----
// INTERFACES
// -----

// DatabasePutter wraps putting key/value pairs in a database.
type DatabasePutter interface {
	Put(key []byte, value []byte) error
}

// DatabaseGetter wraps getting and testing keys in a database.
type DatabaseGetter interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
}

// DatabaseBatch is a write-only database that commits changes to its host database
// when Write is called.
type DatabaseBatch interface {
	DatabasePutter
	ValueSize() int
	Write() error
}

// Database wraps all database operations.
type Database interface {
	DatabasePutter
	DatabaseGetter
	Delete(key []byte) error
	NewBatch() DatabaseBatch
	Close()
}

// Trie is a Ethereum Merkle Trie.
type Trie interface {
	TryGet(key []byte) ([]byte, error)
	TryUpdate(key, value []byte) error
	TryDelete(key []byte) error
	CommitTo(trie.DatabaseWriter) (common.Hash, error)
	Hash() common.Hash
	NodeIterator(startKey []byte) trie.NodeIterator
	GetKey([]byte) []byte // TODO: Remove when SecureTrie is removed.
}

// -----
// TRIE DATABASE
// -----

// TrieDatabase wraps access to tries and contract code.
type TrieDatabase interface {
	// OpenTrie opens the main account trie.
	OpenTrie(root common.Hash) (Trie, error)
	// OpenStorageTrie opens the storage trie of an account.
	OpenStorageTrie(addrHash, root common.Hash) (Trie, error)
	ContractCode(addrHash, codeHash common.Hash) ([]byte, error)
	ContractCodeSize(addrHash, codeHash common.Hash) (int, error)
	// CopyTrie returns an independent copy of the given trie.
	CopyTrie(t Trie) Trie
}

// trieDatabase implements TrieDatabase.
type trieDatabase struct {
	mu            sync.Mutex
	db            Database
	pastTries     []*trie.SecureTrie
	codeSizeCache *lru.Cache
}

// NewTrieDatabase creates a backing store for statements. The returned
// database is safe for concurrent use and retains cached trie nodes
// in memory.
func NewTrieDatabase(db Database) TrieDatabase {
	csc, _ := lru.New(codeSizeCacheSize)
	return &trieDatabase{
		db:            db,
		codeSizeCache: csc,
	}
}

// OpenTrie implements TrieDatabase.
func (tdb *trieDatabase) OpenTrie(root common.Hash) (Trie, error) {
	tdb.mu.Lock()
	defer tdb.mu.Unlock()
	for i := len(tdb.pastTries) - 1; i >= 0; i-- {
		if tdb.pastTries[i].Hash() == root {
			return cachedTrie{tdb.pastTries[i].Copy(), tdb}, nil
		}
	}
	tr, err := trie.NewSecure(root, tdb.db, MaxTrieCacheGen)
	if err != nil {
		return nil, err
	}
	return cachedTrie{tr, tdb}, nil
}

// OpenStorageTrie implements TrieDatabase.
func (tdb *trieDatabase) OpenStorageTrie(addrHash, root common.Hash) (Trie, error) {
	return trie.NewSecure(root, tdb.db, 0)
}

// ContractCode implements TrieDatabase.
func (tdb *trieDatabase) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
	code, err := tdb.db.Get(codeHash[:])
	if err == nil {
		tdb.codeSizeCache.Add(codeHash, len(code))
	}
	return code, err
}

// ContractCodeSize implements TrieDatabase.
func (tdb *trieDatabase) ContractCodeSize(addrHash, codeHash common.Hash) (int, error) {
	if cached, ok := tdb.codeSizeCache.Get(codeHash); ok {
		return cached.(int), nil
	}
	code, err := tdb.ContractCode(addrHash, codeHash)
	if err == nil {
		tdb.codeSizeCache.Add(codeHash, len(code))
	}
	return len(code), err
}

// CopyTrie implements TrieDatabase.
func (tdb *trieDatabase) CopyTrie(t Trie) Trie {
	switch t := t.(type) {
	case cachedTrie:
		return cachedTrie{t.SecureTrie.Copy(), tdb}
	case *trie.SecureTrie:
		return t.Copy()
	default:
		panic(fmt.Errorf("unknown trie type %T", t))
	}
}

// pushTries add the passed tries to it's list of past tries.
func (tdb *trieDatabase) pushTrie(t *trie.SecureTrie) {
	tdb.mu.Lock()
	defer tdb.mu.Unlock()
	if len(tdb.pastTries) >= maxPastTries {
		copy(tdb.pastTries, tdb.pastTries[1:])
		tdb.pastTries[len(tdb.pastTries)-1] = t
	} else {
		tdb.pastTries = append(tdb.pastTries, t)
	}
}

// -----
// CACHED TRIE
// -----

// cachedTrie inserts its trie into a cachingDB on commit.
type cachedTrie struct {
	*trie.SecureTrie
	db *trieDatabase
}

// CommitTo writes the cached trie into the passed trie writer.
func (ct cachedTrie) CommitTo(dbw trie.DatabaseWriter) (common.Hash, error) {
	root, err := ct.SecureTrie.CommitTo(dbw)
	if err == nil {
		ct.db.pushTrie(ct.SecureTrie)
	}
	return root, err
}

// -----
// STATEMENTS DATABASE
// -----

// StatementsDB persists statements and organises them in a trie.
type StatementsDB struct {
	db   TrieDatabase
	trie Trie
}

// Create a new Statements database from a given trie.
func New(root common.Hash, db TrieDatabase) (*StatementsDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	return &StatementsDB{
		db:   db,
		trie: tr,
	}, nil
}
