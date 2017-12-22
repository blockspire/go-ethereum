package observer

import "github.com/ethereum/go-ethereum/common"

// ObsBlock represents one block on the observer chain
// Signature is based on the hash of the RLP encoding of the struct while the "Signature" field is set to nil.
type ObsBlock struct {
	PrevHash      common.Hash
	Number        uint64
	UnixTime      uint64
	TrieRoot      common.Hash // root hash of a trie.Trie structure that is updated for every new block
	SignatureType string      // "ECDSA"
	Signature     []byte      // 65-byte ECDSA signature
}
