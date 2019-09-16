package jahcoin

import (
	"time"
)

type Blockchain struct {
	GenesisBlock *Block
}

// NewBlockchain returns a pointer to a blockchain and any errors
func NewBlockchain(transactionsPerBlock int) (*Blockchain, error) {
	// TODO: check for valid transactionsPerBlock

	b := &Blockchain{
		GenesisBlock: &Block{},
	}

	return b, nil
}

// Block is a node of a blockchain
type Block struct {
	PrevHash     string
	Hash         string
	Timestamp    time.Time
	transactions merkele
}

// Transaction is a transaction between 2 parties
type Transaction struct {
	Amount    float64
	Recipient string
	Sender    string
}
