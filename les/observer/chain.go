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
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
)

// ErrNoFirstBlock - ...
var ErrNoFirstBlock = errors.New("First block not found in observer chain")

// Chain ...
type Chain struct {
	config      *params.ChainConfig // Do we need any configuration?
	chainDb     ethdb.Database
	firstBlock  *Block
	currentBlok *Block
	PrivateKey  *ecdsa.PrivateKey
}

// Block returns a single block by its
func (o *Chain) Block(number uint64) (*Block, error) {
	b := GetBlock(o.chainDb, number)
	return b, nil
}

// LockAndGetTrie lock trie mutex and get r/w access to the current observer trie
func (o *Chain) LockAndGetTrie() *trie.Trie {
	return &trie.Trie{}
}

// UnlockTrie unlock trie mutex
func (o *Chain) UnlockTrie() {

}

// CreateBlock commits current trie and seals a new block; continues using the same trie
// values are persistent, we will care about garbage collection later
func (o *Chain) CreateBlock() *Block {
	return &Block{}
}

// AutoCreateBlocks ...
// creates a new block periodically until chain is closed; non-blocking, starts a goroutine
func (o *Chain) AutoCreateBlocks(period time.Duration) {

}

// Close closes the chain
func (o *Chain) Close() {

}

// NewChain returns a fully initialised Observer chain
// using information available in the database
func NewChain(db ethdb.Database) (*Chain, error) {
	oc := &Chain{
		chainDb: db,
	}
	oc.firstBlock, _ = oc.Block(0)
	if oc.firstBlock == nil {
		// We have to create the first block
		header := Header{Number: 0}
		b := NewBlockWithHeader(&header)
		oc.firstBlock = b
		oc.currentBlok = b
		err := WriteBlock(db, oc.firstBlock)
		if err != nil {
			return oc, err
		}
	}
	return oc, nil
}
