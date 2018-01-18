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
	"testing"

	"github.com/ethereum/go-ethereum/les/observer"
)

// -----
// TESTS
// -----

// TestStatement tests creating and accessing statements.
func TestStatement(t *testing.T) {
	k := []byte("foo")
	v := []byte("bar")
	st := observer.NewStatement(k, v)
	if tk := st.Key(); !bytes.Equal(tk, k) {
		t.Errorf("returned key %v is not value %v", tk, k)
	}
	if tv := st.Value(); !bytes.Equal(tv, v) {
		t.Errorf("returned value %v is not value %v", tv, v)
	}
	var b bytes.Buffer
	err := st.EncodeRLP(&b)
	if err != nil {
		t.Errorf("encoding to RLP returned error: %v", err)
	}
	if b.Len() != 12 {
		t.Errorf("buffer len is %v", b.Len())
	}
}
