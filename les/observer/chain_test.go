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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/les/observer"
)

func TestNewChainHasFistBlockWithNumberZero(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()
	//testdb, _ := ethdb.NewLDBDatabase("./xxx", 10, 256)

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("generation of private key failed")
	}

	c, err := observer.NewChain(testdb, privKey)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}
	if c.FirstBlock().Number().Uint64() != 0 {
		t.Errorf("First block number is not zero")
	}
	if c.CurrentBlock().Number().Uint64() != 0 {
		t.Errorf("Last block number is not zero")
	}
}

func TestWeCanRetrieveFirstBlockFromNewChain(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("generation of private key failed")
	}

	c, err := observer.NewChain(testdb, privKey)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}

	fBlock, err := c.Block(0)
	if err != nil {
		t.Errorf("Retrieve block error = %v", err)
	}
	if fBlock.Number().Uint64() != 0 {
		t.Errorf("First Block has no zero number")
	}
}

func TestCanPersistSecondBlock(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("generation of private key failed")
	}

	c, err := observer.NewChain(testdb, privKey)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}

	secondBlock := observer.NewBlock(privKey)
	if err := observer.WriteBlock(testdb, secondBlock); err != nil {
		t.Errorf("WriteBlock error = %v", err)
	}

	b2 := c.FirstBlock().CreateSuccessor(common.Hash{}, privKey)
	observer.WriteBlock(testdb, b2)

	b2Retrieved, err := c.Block(1)
	if err != nil {
		t.Errorf("Retrieve block error = %v", err)
	}
	if b2Retrieved.Number().Uint64() != 1 {
		t.Errorf("Second Block Number is not 1")
	}
}

func TestWeCanLockAndGetTrieOnNewChain(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("generation of private key failed")
	}

	c, err := observer.NewChain(testdb, privKey)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}

	// Locking and getting the statement trie
	observerTrie, err := c.LockAndGetTrie()
	if err != nil {
		t.Error("Can not unlock trie on new chain")
	}

	err = observerTrie.TryUpdate([]byte("SomeKey"), []byte("SomeValue"))
	if err != nil {
		t.Errorf("updating trie failed")
	}
}

func TestWeCanLockAndGetTrieOnce(t *testing.T) {
	testdb, _ := ethdb.NewMemDatabase()

	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("generation of private key failed")
	}

	c, err := observer.NewChain(testdb, privKey)
	if err != nil {
		t.Errorf("NewChain() error = %v", err)
		return
	}

	// Locking and getting the statement trie
	observerTrie, err := c.LockAndGetTrie()
	if err != nil {
		t.Error("New chain has no trie root")
	}
	if observerTrie.Root() == nil {
		t.Error("Non nil trie has no Root()")
	}

	observerTrie2, err := c.LockAndGetTrie()
	if err == nil {
		t.Error("Locked trie happened to be locked twice :(")
	}
	if err != observer.ErrTrieIsAlreadyLocked {
		t.Error("Error type is not what expected")
	}
	if observerTrie2 != nil {
		t.Error("Observer trie should be nil")
	}
}
