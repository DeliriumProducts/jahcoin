package jahcoin

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"math"
	"math/rand"
	"sync"
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
	// Gekyume is the reciever of the first transaction in the genesis block
	Gekyume                  ed25519.PublicKey
	InitialTransactionAmount float64
	Hasher                   hash.Hash
}

// Blockchain contains all of the blocks and the configuration options
type Blockchain struct {
	GekyumeBlock *Block
	m            *sync.Mutex
	CurrentBlock *Block
	Config       Config
}

// Block is a node of a blockchain
type Block struct {
	Prev         *Block
	PrevHash     []byte
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

// Hash returns the hash of the block
func (b *Block) Hash() [sha256.Size]byte {
	// this is slow, replace with smth else
	// TODO: somehow omit the prev pointer,
	// as they are not reproducible
	// (say you re-run the program, the pointers would be different)
	return sha256.Sum256([]byte(fmt.Sprintf("%v", b)))
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

	t := Transaction{
		Amount:    c.InitialTransactionAmount,
		Recipient: c.Gekyume,
	}

	gekyumeBlock := &Block{
		Prev:         nil,
		PrevHash:     nil,
		Transactions: []Transaction{t},
	}

	b := &Blockchain{
		GekyumeBlock: gekyumeBlock,
		CurrentBlock: gekyumeBlock,
		m:            &sync.Mutex{},
		Config:       *c,
	}

	b.Mine()

	return b, nil
}

func (b *Blockchain) Mine() {
	for {
		hash := b.CurrentBlock.Hash()
		fmt.Println(hex.EncodeToString(hash[:]))
		b.CurrentBlock.Nonce = rand.Int()
		fmt.Println(b.CurrentBlock.Nonce)
	}
}

func (b *Blockchain) AddTransaction(t *Transaction) error {
	b.m.Lock()
	defer b.m.Unlock()

	if len(b.CurrentBlock.Transactions) < b.Config.TransactionsPerBlock {
		b.CurrentBlock.Transactions = append(b.CurrentBlock.Transactions, *t)
	}

	return nil
}
