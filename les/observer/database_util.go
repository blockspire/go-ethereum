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
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

var observerBlockHashPrefix = []byte("o") // observerBlockHashPrefix + hash -> num (uint64 big endian)

// WriteBlock serializes and writes block into the database
// TODO: make it work :)
func WriteBlock(db ethdb.Putter, block *Block) error {
	data, err := rlp.EncodeToBytes(block)
	if err != nil {
		return err
	}

	hash := block.Hash().Bytes()
	num := block.Number().Uint64()
	encNum := encodeBlockNumber(num)
	key := append(observerBlockHashPrefix, hash...)

	if err := db.Put(key, data); err != nil {
		log.Crit("Failed to store header", "err", err)
	}
	return nil
}
