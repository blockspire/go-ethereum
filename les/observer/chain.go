package observer

import (
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

// Chain ...
type Chain struct {
}

// NewChain ...
func NewChain(db ethdb.Database) *Chain {
	return &Chain{}
}

// GetHead ...
func (o *Chain) GetHead() *Block {
	return &Block{}
}

// GetBlock ...
func (o *Chain) GetBlock(index uint64) *Block {
	return &Block{}
}

// LockAndGetTrie lock trie mutex and get r/w access to the current observer trie
func (o *Chain) LockAndGetTrie() *trie.Trie {
	return &trie.Trie{}
}

// UnlockTrie unlock trie mutex
func (o *Chain) UnlockTrie() {

}

// CreateBlock commits current trie and seals a new block; continues using the same trie
// values are persistent, we will care about garbage collection later
func (o *Chain) CreateBlock() *Block {
	return &Block{}
}

// AutoCreateBlocks ...
// creates a new block periodically until chain is closed; non-blocking, starts a goroutine
func (o *Chain) AutoCreateBlocks(period time.Duration) {

}

// Close closes the chain
func (o *Chain) Close() {

}
