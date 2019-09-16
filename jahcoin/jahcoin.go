package jahcoin

import (
	"crypto/ed25519"
	"errors"
	"math"
	"time"
)

var (
	ErrInvalidTransactionsPerBlock = errors.New("jahcoin: TransactionsPerBlock must be a proper logarithm of base 2")
	ErrInvalidDifficulty           = errors.New("jahcoin: Difficulty must be a number higher than 0")
)

// Config contains configuration options for a blockchain
type Config struct {
	TransactionsPerBlock int
	Difficulty           int
}

type Blockchain struct {
	GekyumeBlock *Block
	// TODO: Make it thread safe!!!!??????@@?@?
	CurrentBlock *Block
	config       Config
	Blocks       []*Block
}

// Block is a node of a blockchain
type Block struct {
	PrevHash     string
	Hash         string
	Timestamp    time.Time
	Transactions []Transaction
	Nonce        int
	JahRoot      []byte
}

// Transaction is a transaction between 2 parties
type Transaction struct {
	Amount    float64
	Recipient ed25519.PublicKey
	Sender    ed25519.PublicKey
}

// NewBlockchain returns a pointer to a blockchain and any errors
func NewBlockchain(c *Config) (*Blockchain, error) {

	temp := math.Log2(float64(c.TransactionsPerBlock))

	if temp != float64(int64(temp)) {
		return nil, ErrInvalidTransactionsPerBlock
	}

	if c.Difficulty < 1 {
		return nil, ErrInvalidDifficulty
	}

	b := &Blockchain{
		GekyumeBlock: &Block{},
		config:       *c,
	}

	b.CurrentBlock = b.GekyumeBlock
	b.Blocks = []*Block{b.CurrentBlock}

	return b, nil
}

// func (b *Blockchain) AddBlock(bk *Block) error {
// 	return nil
// }

func (b *Blockchain) AddTransaction(t *Transaction) error {
	if len(b.CurrentBlock.Transactions) < b.config.TransactionsPerBlock {
		b.CurrentBlock.Transactions = append(b.CurrentBlock.Transactions, *t)
	}

	if len(b.CurrentBlock.Transactions) == b.config.TransactionsPerBlock {
		b.Blocks = append(b.Blocks, b.CurrentBlock)
	}

	return nil

}
