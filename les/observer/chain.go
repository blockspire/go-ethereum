package observer

import (
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

// ObsChain ...
type ObsChain struct {
}

// NewObsChain ...
func NewObsChain(db ethdb.Database) *ObsChain {
	return &ObsChain{}
}

// GetHead ...
func (o *ObsChain) GetHead() *ObsBlock {
	return &ObsBlock{}
}

// GetBlock ...
func (o *ObsChain) GetBlock(index uint64) *ObsBlock {
	return &ObsBlock{}
}

// LockAndGetTrie lock trie mutex and get r/w access to the current observer trie
func (o *ObsChain) LockAndGetTrie() *trie.Trie {
	return &trie.Trie{}
}

// UnlockTrie unlock trie mutex
func (o *ObsChain) UnlockTrie() {

}

// CreateBlock commits current trie and seals a new block; continues using the same trie
// values are persistent, we will care about garbage collection later
func (o *ObsChain) CreateBlock() *ObsBlock {
	return &ObsBlock{}
}

// AutoCreateBlocks ...
// creates a new block periodically until chain is closed; non-blocking, starts a goroutine
func (o *ObsChain) AutoCreateBlocks(period time.Duration) {

}

// Close closes the chain
func (o *ObsChain) Close() {

}
