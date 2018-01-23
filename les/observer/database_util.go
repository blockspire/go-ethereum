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
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	observerPrefix = []byte("obs-")      // observerBlockHashPrefix + hash -> num (uint64 big endian)
	lastBlockKey   = []byte("LastBlock") // keeps track of the last observer block
)

// GetBlock retrieves an entire block corresponding to the number, assembling it
// back from the stored header (and statements?). If either the header or body could
// not be retrieved nil is returned.
func GetBlock(db trie.DatabaseReader, number uint64) *Block {
	data := GetBlockRLP(db, number)
	if len(data) == 0 {
		return nil
	}
	b := new(Block)
	if err := rlp.Decode(bytes.NewReader(data), b); err != nil {
		log.Error("Invalid block RLP", "number", number, "err", err)
		return nil
	}
	return b
}

// GetBlockRLP retrieves a block in its raw RLP database encoding, or nil
// if the header's not found.
func GetBlockRLP(db trie.DatabaseReader, number uint64) rlp.RawValue {
	data, _ := db.Get(observerKey(number))
	return data
}

// WriteBlock serializes and writes block into the database
func WriteBlock(db ethdb.Putter, block *Block) error {
	var buf bytes.Buffer
	err := block.EncodeRLP(&buf)
	if err != nil {
		return err
	}
	hash := block.Hash().Bytes()
	key := append(observerPrefix, hash...)
	if err := db.Put(key, buf.Bytes()); err != nil {
		log.Crit("Failed to store observer block data", "err", err)
	}
	return nil
}

// WriteHeadObserverBlockHash writes last block hash to DB under key headBlockKey
func WriteHeadObserverBlockHash(db ethdb.Putter, hash common.Hash) error {
	if err := db.Put(lastBlockKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last observer block's hash", "err", err)
	}
	return nil
}

// -----
// HELPER
// -----

// observerKey calculates the observer key for a given block number.
func observerKey(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return append(observerPrefix, enc...)
}
