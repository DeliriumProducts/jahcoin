package jahcoin

import "crypto/ed25519"

// Transaction is a transaction between 2 parties
type Transaction struct {
	Amount    float64
	Recipient ed25519.PublicKey
	Sender    ed25519.PublicKey
}
