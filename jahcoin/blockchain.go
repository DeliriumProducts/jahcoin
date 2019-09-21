package jahcoin

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
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
	// GekyumeBlock (also known as the GenesisBlock) is the first block in the entire blockchain
	GekyumeBlock *Block
	m            *sync.Mutex
	CurrentBlock *Block
	Config       Config
}

// NewBlockchain returns a pointer to a blockchain and any errors
func NewBlockchain(c *Config) (*Blockchain, error) {
	// TODO: store the entire blockchain on disk using a file or db
	// and parse / read it here.
	// Also have to handle pointing blocks to the previous ones

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
		Timestamp:    time.Now(),
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
			b.CurrentBlock.Hashed = hash[:]
			log.Printf("Nonce found! %v, hash is: %v", b.CurrentBlock.Nonce, hx)
			break
		}

		b.CurrentBlock.Nonce = rand.Int()
	}

	// Make the CurrentBlock point to the new one
	minedBlock := *b.CurrentBlock
	b.CurrentBlock = &Block{
		prev:         &minedBlock,
		PrevHash:     minedBlock.Hashed,
		Timestamp:    time.Now(),
		Transactions: []Transaction{},
	}
}

func (b *Blockchain) AddTransaction(t *Transaction) {
	// TODO: check if the sender has enough cash to send
	// BTC has an indexed db stored with the
	// balances of all public keys, they are collected
	// by traversing the entire blockchain

	b.m.Lock()
	defer b.m.Unlock()

	lx := len(b.CurrentBlock.Transactions)
	if lx < b.Config.TransactionsPerBlock {
		b.CurrentBlock.Transactions = append(b.CurrentBlock.Transactions, *t)

		// If we haven't reached the limit yet, wait for more transactions
		if lx+1 < b.Config.TransactionsPerBlock {
			return
		}
	}

	// Calculate JahRoot (MerkleTree)
	b.CurrentBlock.JahRoot = b.hashTransactions()

	log.Printf("%+v", b.CurrentBlock)
}

// hashTransactions returns the Merkle root / hash
// of all transactions in the current block by
// recursively hashing all elements.
// Example:
// h() being an arbitrary hashing function,`TransactionsPerBlock` set to 8
//	   [A, B, C, D, E, F, G, H] ->

//	   [h(A), h(B), h(C), h(D), h(E), h(F), h(G), h(H)] ->

//	   [h(
//		 h(A) + h(B)
//		),
//		h(
//		 h(C) + h(D)
//		),
//		h(
//		 h(E) + h(F)
//	   ),
//		h(
//		 h(G) + h(H)
//	   )] ->

//	   [h(
//			h(h(A) + h(B)) +
//			h(h(C) + h(D))
//      ),
//		h(
//			h(h(E) + h(F)) +
//			h(h(G) + h(H))
//	   )] ->
//
//	   [h(
//			h(
//				h(h(A) + h(B)) +
//				h(h(C) + h(D))
// 			) +
//			h(
//				h(h(E) + h(F)) +
//				h(h(G) + h(H))
//			)
//	   )] ->
//     e68fe78e064700fe6b98e47dc0758a4f966bd027299b685642c607ea376b7d47
func (b *Blockchain) hashTransactions() []byte {
	// TODO: somehow make concurrent (?)
	// I don't know if that's possible

	// TODO: are we vulnerable to the
	// Second Preimage attack?
	// https://www.wikiwand.com/en/Merkle_tree

	var helper func(levels [][]byte) []byte

	// Helper takes in the current level of the merkle tree
	// and pairs up all elements inside
	helper = func(levels [][]byte) []byte {
		lx := len(levels)
		// The bottom of our recursion:
		// (we have only 2 elements left, and we have to calculate the final root hash)
		// [A, B]
		if lx == 2 {
			return hashVariadic(levels[0], levels[1])
		}

		// the next level is always going to be twice as small
		// [A, B, C, D] ->
		// [AB, CD] ->
		// [ABCD]

		newLevels := make([][]byte, lx/2)

		// keep track of the index in newLevels (has to be seperate value, as it's twice as small
		// maybe there's a  better way to compute the idx without needing another variable?
		idx := 0

		// start pairing each element with its neighbor
		// [A, B, C, D] ->
		// A and B, C and D
		for i := 0; i < len(levels); i += 2 {
			newLevels[idx] = hashVariadic(levels[i], levels[i+1])
			idx++
		}

		return helper(newLevels)
	}

	// get the initial array of hashes (the transactions)
	levels := [][]byte{}

	for _, tx := range b.CurrentBlock.Transactions {
		h, err := tx.Hash()
		if err != nil {
			log.Println(err)
		}

		levels = append(levels, h[:])
	}

	return helper(levels)
}

// hashVariadic takes in n amount of
// elements and returns their sha256 summed hash
func hashVariadic(b ...[]byte) []byte {
	h := sha256.New()

	for _, bt := range b {
		h.Write(bt)
	}

	return h.Sum([]byte{})
}
