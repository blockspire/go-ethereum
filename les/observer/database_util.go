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
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	blockPrefix      = []byte("obs-")      // blockPrefix + num -> Block
	stmtLookupPrefix = []byte("obssl-")    // stmtLookupPrefix + key -> StmtLookupEntry
	lastBlockKey     = []byte("lastBlock") // keeps track of the last observer block
)

// StmtLookupEntry is a positional metadata to help looking up the statement
// inside its block.
type StmtLookupEntry struct {
	BlockNumber uint64
	Index       uint64
}

// GetBlock retrieves an entire block corresponding to the number.
func GetBlock(db DatabaseGetter, number uint64) *Block {
	data, _ := db.Get(mkBlockKey(number))
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

// GetStmtLookupEntry retrieves block number and index of a statement.
func GetStmtLookupEntry(db DatabaseGetter, key []byte) (uint64, uint64, bool) {
	// Retrieve lookup entry.
	data, _ := db.Get(mkStmtLookupKey(key))
	if len(data) == 0 {
		log.Error("Cannot find statement lookup", "key", key)
		return 0, 0, false
	}
	// Decode it.
	var entry StmtLookupEntry
	if err := rlp.DecodeBytes(data, &entry); err != nil {
		log.Error("Invalid lookup entry RLP", "key", key, "err", err)
		return 0, 0, false
	}
	return entry.BlockNumber, entry.Index, true
}

// GetStatement retrieves a specific statement from the database by key. It
// also returns the number of the block and the index of the statement inside
// of it.
func GetStatement(db DatabaseGetter, key []byte) (*Statement, uint64, uint64) {
	// Retrieve block number and statement index.
	blockNumber, stmtIndex, ok := GetStmtLookupEntry(db, key)
	if !ok {
		return nil, 0, 0
	}
	// Retrieve the block and statement.
	if block := GetBlock(db, blockNumber); block != nil {
		if stmt := block.StatementByIndex(stmtIndex); stmt != nil {
			return stmt, blockNumber, stmtIndex
		}
	}
	// Not found.
	log.Error("Statement referenced missing", "block number", blockNumber, "index", stmtIndex)
	return nil, 0, 0
}

// WriteBlock serializes and writes block into the database
func WriteBlock(db DatabasePutter, block *Block) error {
	var buf bytes.Buffer
	err := block.EncodeRLP(&buf)
	if err != nil {
		return err
	}
	if err := db.Put(mkBlockKey(block.header.Number), buf.Bytes()); err != nil {
		log.Crit("Failed to store observer block data", "err", err)
	}
	return nil
}

// WriteLastObserverBlockHash writes last block hash to DB under key headBlockKey
func WriteLastObserverBlockHash(db DatabasePutter, hash common.Hash) error {
	if err := db.Put(lastBlockKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last observer block's hash", "err", err)
	}
	return nil
}

// -----
// HELPER
// -----

// mkBlockKey creates the database key for a given block number.
// Ex: obs-0, obs-124
func mkBlockKey(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return append(blockPrefix, enc...)
}

// mkStmtLookupKey creates the database key for a given statement lookup key.
// Ex: obssl-foo, obssl-bar
func mkStmtLookupKey(key []byte) []byte {
	return append(stmtLookupPrefix, key...)
}
