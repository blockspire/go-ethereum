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
	"crypto/ecdsa"
	"errors"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

// ErrNoFirstBlock - ...
var ErrNoFirstBlock = errors.New("First block not found in observer chain")

// ErrNoBlock if we can not retrieve requested block
var ErrNoBlock = errors.New("Block not found in observer chain")

// ErrTrieIsAlreadyLocked if trie is locked already
var ErrTrieIsAlreadyLocked = errors.New("Can not unlock, Observer trie is already locked, sorry")

const ( // statuses for statement trie
	locked = iota
	unlocked
	unlocking
)

// -----
// CHAIN
// -----

// Chain represents the canonical observer chain given a database with a
// genesis block.
type Chain struct {
	db           ethdb.Database
	firstBlock   *Block
	currentBlock *Block
	privateKey   *ecdsa.PrivateKey
	trieStatus   atomic.Value // Stores the statement trie locked status ( locked/unlocked/unlocking )
}

// NewChain returns a fully initialised Observer chain
// using information available in the database
func NewChain(db ethdb.Database, privKey *ecdsa.PrivateKey) (*Chain, error) {
	oc := &Chain{
		db:         db,
		privateKey: privKey,
	}
	oc.trieStatus.Store(unlocked)
	firstBlock := GetBlock(db, 0)
	if firstBlock == nil {
		firstBlock = NewBlock(privKey)
	}
	oc.firstBlock = firstBlock
	oc.currentBlock = firstBlock
	if err := WriteBlock(db, firstBlock); err != nil {
		return nil, err
	}
	if WriteLastObserverBlockHash(db, firstBlock.Hash()) != nil {
		return nil, nil
	}
	return oc, nil
}

// Block returns a single block by its
func (o *Chain) Block(number uint64) (*Block, error) {
	b := GetBlock(o.db, number)
	if b == nil {
		return nil, ErrNoBlock
	}
	return b, nil
}

// FirstBlock returns Observer Chain's first block, aka. Genesis block.
// We are not calling it genesis block because will be different from node to node.
func (o *Chain) FirstBlock() *Block {
	return o.firstBlock
}

// CurrentBlock returns the current last block on the chain
func (o *Chain) CurrentBlock() *Block {
	return o.currentBlock
}

// LockAndGetTrie lock trie mutex and get r/w access to the current observer trie
func (o *Chain) LockAndGetTrie() (*trie.Trie, error) {
	if sts := o.trieStatus.Load(); sts == nil || sts == unlocked {
		o.trieStatus.Store(locked)
		tr, err := trie.New(o.currentBlock.TrieRoot(), trie.NewDatabase(o.db))
		if err == nil {
			return tr, nil
		}
	}
	return nil, ErrTrieIsAlreadyLocked
}

// UnlockTrie unlock trie mutex
func (o *Chain) UnlockTrie() {
	// check if trie is locked
	// if locked, commit trie, save block, then unlock trie
	if sts := o.trieStatus.Load(); sts == locked {

	}

}

// CreateBlock commits current trie and seals a new block; continues using the same trie
// values are persistent, we will care about garbage collection later
func (o *Chain) CreateBlock() *Block {
	t, err := o.LockAndGetTrie()
	if err == nil {
		t.Commit(nil)
		return o.CurrentBlock().CreateSuccessor(o.CurrentBlock().TrieRoot(), o.privateKey)
	}
	return o.CurrentBlock().CreateSuccessor(o.CurrentBlock().TrieRoot(), o.privateKey)
}

// AutoCreateBlocks ...
// creates a new block periodically until chain is closed; non-blocking, starts a goroutine
func (o *Chain) AutoCreateBlocks(period time.Duration) {

}

// Close closes the chain
func (o *Chain) Close() {

}
