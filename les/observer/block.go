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
	"encoding/binary"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// -----
// BLOCK DATA
// -----

// blockData contains the data fields of a block opposite to
// internal data. Signature is based on the hash of the RLP encoding
// of the struct while the Signature field is set to nil.
type blockData struct {
	PrevHash      common.Hash `json:"prevHash"      gencodec:"required"`
	Number        uint64      `json:"number"        gencodec:"required"`
	UnixTime      uint64      `json:"unixTime"      gencodec:"required"`
	Statements    common.Hash `json:"statements"    gencodec:"required"`
	SignatureType string      `json:"signatureType" gencodec:"required"`
	Signature     []byte      `json:"signature"     gencodec:"required"`
}

// hash returns the block hash of the header, which is simply the keccak256
// hash of its RLP encoding.
func (d *blockData) hash() common.Hash {
	return rlpHash(d)
}

// sign adds a signature to the block data by the given private key.
func (d *blockData) sign(privKey *ecdsa.PrivateKey) {
	unsignedData := &blockData{
		PrevHash:      d.PrevHash,
		Number:        d.Number,
		UnixTime:      d.UnixTime,
		Statements:    d.Statements,
		SignatureType: d.SignatureType,
	}
	rlp, _ := rlp.EncodeToBytes(unsignedData)
	d.Signature, _ = crypto.Sign(crypto.Keccak256(rlp), privKey)
}

// -----
// BLOCK
// -----

// Block represents one block on the observer chain.
type Block struct {
	data *blockData

	// Caches.
	hash atomic.Value
	size atomic.Value
}

// NewBlock creates a new block.
// TODO: More details about arguments.
func NewBlock(sts []*Statement, privKey *ecdsa.PrivateKey) *Block {
	b := &Block{
		data: &blockData{
			PrevHash:      common.Hash{},
			Number:        0,
			UnixTime:      uint64(time.Now().Unix()),
			SignatureType: "ECDSA",
		},
	}
	if len(sts) == 0 {
		b.data.Statements = types.EmptyRootHash
	} else {
		b.data.Statements = types.DeriveSha(Statements(sts))
	}
	b.data.sign(privKey)
	return b
}

// Number returns the block number as big.Int.
func (b *Block) Number() *big.Int {
	return new(big.Int).SetUint64(b.data.Number)
}

// EncodedNumber returns the block number in a big endian
// encoded way.
func (b *Block) EncodedNumber() []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, b.data.Number)
	return enc
}

// Statements returns the hash of the block statements.
func (b *Block) Statements() common.Hash {
	return b.data.Statements
}

// Hash returns the keccak256 hash of the block's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.data.hash()
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
