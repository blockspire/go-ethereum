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
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// -----
// HEADER
// -----

// Header contains the header fields of a block opposite to
// internal data.
type Header struct {
	PrevHash      common.Hash `json:"prevHash"      gencodec:"required"`
	Number        uint64      `json:"number"        gencodec:"required"`
	UnixTime      uint64      `json:"unixTime"      gencodec:"required"`
	Statements    common.Hash `json:"statements"    gencodec:"required"`
	SignatureType string      `json:"signatureType" gencodec:"required"`
	Signature     []byte      `json:"signature"     gencodec:"required"`
}

// Hash returns the block hash of the header, which is simply the keccak256
// hash of its RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h)
}

// -----
// BLOCK
// -----

// Block represents one block on the observer chain.
// Signature is based on the hash of the RLP encoding of
// the struct while the "Signature" field is set to nil.
type Block struct {
	header *Header

	// Caches.
	hash atomic.Value
	size atomic.Value
}

// NewBlock creates a new block.
// TODO: More details about arguments.
func NewBlock(txs []*Statement, privKey *ecdsa.PrivateKey) *Block {
	b := &Block{
		header: &Header{
			PrevHash:      common.Hash{},
			Number:        0,
			UnixTime:      uint64(time.Now().Unix()),
			SignatureType: "ECDSA",
		},
	}
	if len(txs) == 0 {
		b.header.Statements = types.EmptyRootHash
	} else {
		b.header.Statements = types.DeriveSha(Statements(txs))
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
			Statements:    b.header.Statements,
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

// Statements returns the hash of the block statements.
func (b *Block) Statements() common.Hash {
	return b.header.Statements
}

// Hash returns the keccak256 hash of the block's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

// Size returns the storage size of the block.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}
