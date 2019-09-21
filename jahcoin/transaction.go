package jahcoin

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/gob"
)

// Transaction is a transaction between 2 parties
type Transaction struct {
	Amount    float64
	Recipient ed25519.PublicKey
	Sender    ed25519.PublicKey
}

func (t *Transaction) Hash() ([sha256.Size]byte, error) {
	bf := &bytes.Buffer{}
	if err := gob.NewEncoder(bf).Encode(t); err != nil {
		return [sha256.Size]byte{}, err
	}

	return sha256.Sum256(bf.Bytes()), nil
}
