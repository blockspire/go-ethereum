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
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

// Header contains the header fields of a block opposite to
// internal data.
type Header struct {
	PrevHash      common.Hash `json:"prevHash"      gencodec:"required"`
	Number        uint64      `json:"number"        gencodec:"required"`
	UnixTime      uint64      `json:"unixTime"      gencodec:"required"`
	TrieRoot      common.Hash `json:"trieRoot"      gencodec:"required"`
	SignatureType string      `json:"signatureType" gencodec:"required"`
	Signature     []byte      `json:"signature"     gencodec:"required"`
}

// Block represents one block on the observer chain.
// Signature is based on the hash of the RLP encoding of
// the struct while the "Signature" field is set to nil.
type Block struct {
	header *Header
}

// NewBlock creates a new block.
// TODO: More details about arguments.
func NewBlock(txs []*types.Transaction, privKey *ecdsa.PrivateKey) *Block {
	b := &Block{
		header: &Header{
			PrevHash:      common.Hash{},
			Number:        0,
			UnixTime:      uint64(time.Now().Unix()),
			SignatureType: "ECDSA",
		},
	}
	if len(txs) == 0 {
		b.header.TrieRoot = types.EmptyRootHash
	} else {
		b.header.TrieRoot = types.DeriveSha(types.Transactions(txs))
	}
	b.Sign(privKey)
	return b
}

// Sign adds a signature to the block by the given private key.
func (b *Block) Sign(privKey *ecdsa.PrivateKey) {
	unsignedBlock := Block{
		header: &Header{
			PrevHash:      b.header.PrevHash,
			Number:        b.header.Number,
			UnixTime:      b.header.UnixTime,
			TrieRoot:      b.header.TrieRoot,
			SignatureType: b.header.SignatureType,
		},
	}
	rlp, _ := rlp.EncodeToBytes(unsignedBlock)
	b.header.Signature, _ = crypto.Sign(crypto.Keccak256(rlp), privKey)
}

// Number returns the block number as big.Int.
func (b *Block) Number() *big.Int {
	return new(big.Int).SetUint64(b.header.Number)
}

// TrieRoot returns the hash of the trie root.
func (b *Block) TrieRoot() common.Hash {
	return b.header.TrieRoot
}

// rlpHash calculates a hash out of the passed data.
func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
