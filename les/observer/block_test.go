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

package observer_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/les/observer"
	"github.com/ethereum/go-ethereum/rlp"
)

// -----
// TESTS
// -----

// TestStatement tests creating and accessing statements.
func TestStatement(t *testing.T) {
	k := []byte("foo")
	v := []byte("bar")
	st := observer.NewStatement(k, v)
	// Testing simple access.
	if tk := st.Key(); !bytes.Equal(tk, k) {
		t.Errorf("returned key %v is not value %v", tk, k)
	}
	if tv := st.Value(); !bytes.Equal(tv, v) {
		t.Errorf("returned value %v is not value %v", tv, v)
	}
	// Testing encoding and decoding.
	var buf bytes.Buffer
	err := st.EncodeRLP(&buf)
	if err != nil {
		t.Errorf("encoding to RLP returned error: %v", err)
	}
	var tst observer.Statement
	err = tst.DecodeRLP(rlp.NewStream(&buf, 0))
	if err != nil {
		t.Errorf("decoding from RLP returned error: %v", err)
	}
	if tk := tst.Key(); !bytes.Equal(tk, k) {
		t.Errorf("returned decoded key %v is not value %v", tk, k)
	}
	if tv := tst.Value(); !bytes.Equal(tv, v) {
		t.Errorf("returned decoded value %v is not value %v", tv, v)
	}
	sthb := st.Hash().Bytes()
	tsthb := st.Hash().Bytes()
	if !bytes.Equal(sthb, tsthb) {
		t.Errorf("hashes of original and encoded/decoded one differ")
	}
}

// TestEmptyBlock tests creating and accessing empty blocks.
func TestEmptyBlock(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Errorf("generation of private key failed")
	}
	b := observer.NewBlock(nil, privKey)
	if b.Number().Uint64() != 0 {
		t.Errorf("number of new block is not 0")
	}
	// Testing encoding and decoding.
	var buf bytes.Buffer
	err = b.EncodeRLP(&buf)
	if err != nil {
		t.Errorf("encoding to RLP returned error: %v", err)
	}
	var tb observer.Block
	err = tb.DecodeRLP(rlp.NewStream(&buf, 0))
	if err != nil {
		t.Errorf("decoding from RLP returned error: %v", err)
	}
	if tb.Number().Uint64() != 0 {
		t.Errorf("number of encoded/decoded block is not 0")
	}
}

// TestStatementsBlock tests creating and accessing blocks
// containing statements
func TestStatementsBlock(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Errorf("generation of private key failed")
	}
	sts := []*observer.Statement{
		observer.NewStatement([]byte("foo"), []byte("123")),
		observer.NewStatement([]byte("bar"), []byte("456")),
		observer.NewStatement([]byte("baz"), []byte("789")),
	}
	b := observer.NewBlock(sts, privKey)
	if b.Number().Uint64() != 0 {
		t.Errorf("number of new block is not 0")
	}
	// Testing encoding and decoding.
	var buf bytes.Buffer
	err = b.EncodeRLP(&buf)
	if err != nil {
		t.Errorf("encoding to RLP returned error: %v", err)
	}
	var tb observer.Block
	err = tb.DecodeRLP(rlp.NewStream(&buf, 0))
	if err != nil {
		t.Errorf("decoding from RLP returned error: %v", err)
	}
	if tb.Number().Uint64() != 0 {
		t.Errorf("number of encoded/decoded block is not 0")
	}
}
