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
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/les/observer"
)

func TestNewChainHasFistBlockWithNumberZero(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()
	//testdb, _ := ethdb.NewLDBDatabase("./xxx", 10, 256)

	c, err := observer.NewChain(testdb)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}

	firstBlock, err := c.Block(uint64(0))
	if err != nil {
		t.Errorf("Retrieve block error = %v", err)
	}
	if firstBlock.Number().Uint64() != 0 {
		t.Errorf("number of new block is not 0")
	}
}

func TestCanPersistBlock(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()

	sts := []*observer.Statement{
		observer.NewStatement([]byte("foo")),
		observer.NewStatement([]byte("bar")),
		observer.NewStatement([]byte("baz")),
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("generation of private key failed")
	}

	firstBlock := observer.NewBlock(sts, privKey)
	if err := observer.WriteBlock(testdb, firstBlock); err != nil {
		t.Errorf("WriteBlock error = %v", err)
	}
}

func TestWeCanRetrievePersisedBlock(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()

	sts := []*observer.Statement{
		observer.NewStatement([]byte("foo")),
		observer.NewStatement([]byte("bar")),
		observer.NewStatement([]byte("baz")),
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("generation of private key failed")
	}

	firstBlock := observer.NewBlock(sts, privKey)
	if err := observer.WriteBlock(testdb, firstBlock); err != nil {
		t.Errorf("WriteBlock error = %v", err)
	}

	c, err := observer.NewChain(testdb)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}

	fBlock, err := c.Block(uint64(0))
	if err != nil {
		t.Errorf("Retrieve block error = %v", err)
	}
	t.Log(fBlock.Number())
}
