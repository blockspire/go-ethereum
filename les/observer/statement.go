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
	"io"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

// -----
// STATEMENT
// -----

// Statement contains a combination of key and value.
type Statement struct {
	kv keyValue

	// Caches.
	hash atomic.Value
	size atomic.Value
}

// keyValue manages the key and value data of a statement.
type keyValue struct {
	Key   []byte `json:"key"   gencodec:"required"`
	Value []byte `json:"value" gencodec:"required"`

	// Signature values.
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// NewStatement creates a standard statement with key and value.
func NewStatement(key, value []byte) *Statement {
	return newStatement(key, value)
}

// newStatement is the private constructor for the different types
// of statements.
func newStatement(key, value []byte) *Statement {
	// Create modifiable copies of key and value.
	if len(key) > 0 {
		key = common.CopyBytes(key)
	}
	if len(value) > 0 {
		value = common.CopyBytes(value)
	}
	kv := keyValue{
		Key:   key,
		Value: value,
		V:     new(big.Int),
		R:     new(big.Int),
		S:     new(big.Int),
	}
	return &Statement{
		kv: kv,
	}
}

// Key returns copy of the statement key.
func (s *Statement) Key() []byte {
	return common.CopyBytes(s.kv.Key)
}

// Value returns copy of the statement value.
func (s *Statement) Value() []byte {
	return common.CopyBytes(s.kv.Value)
}

// EncodeRLP implements rlp.Encoder.
func (s *Statement) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &s.kv)
}

// DecodeRLP implements rlp.Decoder.
func (s *Statement) DecodeRLP(str *rlp.Stream) error {
	_, size, _ := str.Kind()
	err := str.Decode(&s.kv)
	if err == nil {
		s.size.Store(common.StorageSize(rlp.ListSize(size)))
	}
	return err
}

// Hash hashes the RLP encoding of the statement.
// It uniquely identifies it.
func (s *Statement) Hash() common.Hash {
	if hash := s.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(s)
	s.hash.Store(v)
	return v
}

// Size returns the storage size of the statement.
func (s *Statement) Size() common.StorageSize {
	if size := s.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &s.kv)
	s.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// -----
// STATEMENTS
// -----

// Statements contains a number of statements as slice.
type Statements []*Statement

// Len implements types.DerivableList and returns the number
// of statements.
func (s Statements) Len() int {
	return len(s)
}

// GetRlp implements types.DerivableList and returns the i'th
// statement of s in RLP encoding.
func (s Statements) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// -----
// HELPERS
// -----

// rlpHash calculates a hash out of the passed data.
func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// writeCounter helps counting the written bytes in total.
type writeCounter common.StorageSize

// Write implements io.Writer and counts the written bytes
// in total.
func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}
