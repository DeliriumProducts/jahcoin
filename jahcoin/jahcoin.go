package jahcoin

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"log"
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
	// Difficulty is the leading zeroes needed for a block to be valid
	Difficulty int
	// Gekyume is the reciever of the first transaction in the genesis block
	Gekyume                  ed25519.PublicKey
	InitialTransactionAmount float64
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
	prev         *Block
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
func (b *Block) Hash() ([sha256.Size]byte, error) {
	bf := &bytes.Buffer{}
	if err := gob.NewEncoder(bf).Encode(b); err != nil {
		return [sha256.Size]byte{}, err
	}

	return sha256.Sum256(bf.Bytes()), nil
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
		prev:         nil,
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

// TODO: make concurrent / workerpool
func (b *Blockchain) Mine() {
	for {
		hash, err := b.CurrentBlock.Hash()
		if err != nil {
			log.Println(err)
		}

		hx := hex.EncodeToString(hash[:])
		zeroes := 0
		for _, cx := range hx {
			if cx != '0' {
				break
			}

			zeroes++
		}

		if zeroes >= b.Config.Difficulty {
			log.Printf("Nonce found! %v, hash is: %v", b.CurrentBlock.Nonce, hx)
			return
		}

		b.CurrentBlock.Nonce = rand.Int()
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
