package jahcoin

import (
	"errors"
	"math"
)

var (
	ErrNotSufficentTransactions = errors.New("merkele: There are not enough transactions to create a merkele tree")
	ErrTooManyTransactions      = errors.New("merkele: There are too many transactions to create a merkele tree")
)

// Node is an element of a Merkele tree
type node struct {
	Index int
	Hash  string
	Left  *node
	Right *node
}

// Merkele tree contains all the transactions in a block
// and their combined hashes
type merkele struct {
	Root         *node
	transactions []Transaction
}

// NewMerkele returns a pointer to a Merkele tree and any errors
func newMerkele(transactions []Transaction, transactionsPerBlock int) (*merkele, error) {
	tAmount := len(transactions)

	if tAmount >= transactionsPerBlock {
		return nil, ErrNotSufficentTransactions
	}

	if tAmount < 1 {
		return nil, ErrTooManyTransactions
	}

	m := &merkele{
		transactions: transactions,
	}

	levels := math.Log2(float64(transactionsPerBlock))

	m.Root = &node{}

	return m, nil
}
