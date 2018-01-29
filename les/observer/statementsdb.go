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
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	lru "github.com/hashicorp/golang-lru"
)

const (
	// Number of codehash->size associations to keep.
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
	CopyTrie(Trie) Trie
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
	return nil, nil
}

// OpenStorageTrie implements TrieDatabase.
func (tdb *trieDatabase) OpenStorageTrie(addrHash, root common.Hash) (Trie, error) {
	return nil, nil
}

// ContractCode implements TrieDatabase.
func (tdb *trieDatabase) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
	return nil, nil
}

// ContractCodeSize implements TrieDatabase.
func (tdb *trieDatabase) ContractCodeSize(addrHash, codeHash common.Hash) (int, error) {
	return 0, nil
}

// CopyTrie implements TrieDatabase.
func (tdb *trieDatabase) CopyTrie(Trie) Trie {
	return nil
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
