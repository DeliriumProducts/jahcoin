package jahcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// Block is a node of a blockchain
type Block struct {
	prev     *Block
	PrevHash []byte
	// Hashed is the hashed value of the entire block (excluding the prev pointer)
	Hashed       []byte
	Timestamp    time.Time
	Transactions []Transaction
	Nonce        int
	// JahRoot is the Merkle / Hash value of all transactions in the current block
	JahRoot []byte
}

// Hash returns the hash of the block
func (b *Block) Hash() (*[sha256.Size]byte, error) {
	bf := &bytes.Buffer{}
	if err := gob.NewEncoder(bf).Encode(b); err != nil {
		return &[sha256.Size]byte{}, err
	}

	h := sha256.Sum256(bf.Bytes())
	return &h, nil
}
