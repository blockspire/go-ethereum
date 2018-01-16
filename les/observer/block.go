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
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

// Block represents one block on the observer chain.
// Signature is based on the hash of the RLP encoding of
// the struct while the "Signature" field is set to nil.
type Block struct {
	PrevHash      common.Hash `json:"prevHash"      gencodec:"required"`
	Number        uint64      `json:"number"        gencodec:"required"`
	UnixTime      uint64      `json:"unixTime"      gencodec:"required"`
	TrieRoot      common.Hash `json:"trieRoot"      gencodec:"required"`
	SignatureType string      `json:"signatureType" gencodec:"required"`
	Signature     []byte      `json:"signature"     gencodec:"required"`
}

// NewBlock creates a new block.
// TODO: More details for arguments.
func NewBlock(privKey *ecdsa.PrivateKey) *Block {
	b := &Block{
		PrevHash:      common.Hash{},
		Number:        0,
		UnixTime:      uint64(time.Now().Unix()),
		SignatureType: "ECDSA",
	}
	b.sign(privKey)
	return b
}

// sign adds a signature to the block by the given privKey
func (b *Block) sign(privKey *ecdsa.PrivateKey) {
	unsignedBlock := Block{
		PrevHash:      b.PrevHash,
		Number:        b.Number,
		UnixTime:      b.UnixTime,
		TrieRoot:      b.TrieRoot,
		SignatureType: b.SignatureType,
	}
	rlp, _ := rlp.EncodeToBytes(unsignedBlock)
	sig, _ := crypto.Sign(crypto.Keccak256(rlp), privKey)
	b.Signature = b.Signature.add("sign", sig)
}

// rlpHash calculates a hash out of the passed data.
func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}