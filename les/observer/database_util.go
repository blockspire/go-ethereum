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
	lookupPrefix   = []byte("obsl-")     // lookup of observer statement keys to block hashes
	lastBlockKey   = []byte("lastBlock") // keeps track of the last observer block
)

// StmtLookupEntry is a positional metadata to help looking up the data content of
// a statement given only its hash.
type StmtLookupEntry struct {
	BlockHash common.Hash
	Index     uint64
}

// GetBlock retrieves an entire block corresponding to the number.
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

// GetBlockByHash retrieves a block by hash.
func GetBlockByHash(db trie.DatabaseReader, hash common.Hash) *Block {
	// TODO: Implement!
	return nil
}

// GetBlockRLP retrieves a block in its raw RLP database encoding, or nil
// if the header's not found.
func GetBlockRLP(db trie.DatabaseReader, number uint64) rlp.RawValue {
	data, _ := db.Get(observerKey(number))
	return data
}

// GetStmtLookupEntry retrieves the positional metadata associated with a
// statement hash to allow retrieving the statement by hash.
func GetStmtLookupEntry(db trie.DatabaseReader, hash common.Hash) (common.Hash, uint64) {
	// Load the positional metadata from disk and bail if it fails.
	data, _ := db.Get(append(lookupPrefix, hash.Bytes()...))
	if len(data) == 0 {
		return common.Hash{}, 0
	}
	// Parse and return the contents of the lookup entry.
	var entry StmtLookupEntry
	if err := rlp.DecodeBytes(data, &entry); err != nil {
		log.Error("Invalid lookup entry RLP", "hash", hash, "err", err)
		return common.Hash{}, 0
	}
	return entry.BlockHash, entry.Index
}

// GetStatement retrieves a specific statement from the database, along with
// its added positional metadata.
func GetStatement(db trie.DatabaseReader, hash common.Hash) (*Statement, common.Hash, uint64) {
	// Retrieve hash of the block
	blockHash, stmtIndex := GetStmtLookupEntry(db, hash)
	if blockHash != (common.Hash{}) {
		block := GetBlockByHash(db, blockHash)
		if block == nil || len(block.statements) <= int(stmtIndex) {
			log.Error("Transaction referenced missing", "hash", blockHash, "index", stmtIndex)
			return nil, common.Hash{}, 0
		}
		return block.statements[stmtIndex], blockHash, stmtIndex
	}
	return nil, common.Hash{}, 0
}

// WriteBlock serializes and writes block into the database
func WriteBlock(db ethdb.Putter, block *Block) error {
	var buf bytes.Buffer
	err := block.EncodeRLP(&buf)
	if err != nil {
		return err
	}
	//key := append(observerPrefix, encodeBlockNumber(block.header.Number)...)
	if err := db.Put(observerKey(block.header.Number), buf.Bytes()); err != nil {
		log.Crit("Failed to store observer block data", "err", err)
	}
	return nil
}

// WriteLastObserverBlockHash writes last block hash to DB under key headBlockKey
func WriteLastObserverBlockHash(db ethdb.Putter, hash common.Hash) error {
	if err := db.Put(lastBlockKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last observer block's hash", "err", err)
	}
	return nil
}

// -----
// HELPER
// -----

// observerKey calculates the observer key for a given block number.
// ex: obs-0, obs-124
func observerKey(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return append(observerPrefix, enc...)
}

func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}
