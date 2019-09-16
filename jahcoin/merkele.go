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
	index int
	hash  string
	left  *node
	right *node
}

// Merkele tree contains all the transactions in a block
// and their combined hashes
type merkele struct {
	root         *node
	transactions []Transaction
}

// NewMerkele returns a pointer to a Merkele tree and any errors
func newMerkele(transactions []Transaction, transactionsPerBlock int) (*merkele, error) {
	tAmount := len(transactions)

	if tAmount >= transactionsPerBlock {
		return nil, ErrTooManyTransactions
	}

	if tAmount < 1 {
		return nil, ErrNotSufficentTransactions
	}

	m := &merkele{
		transactions: transactions,
	}

	levels := math.Log2(float64(transactionsPerBlock))

	m.root = &node{}

	var helper func(n *node, currentLevel int)

	helper = func(n *node, currentLevel int) {

		if currentLevel < 1 {
			return
		}

		n.right = &node{}
		helper(n.right, currentLevel-1)

		n.left = &node{}
		helper(n.left, currentLevel-1)
	}

	helper(m.root, int(levels))

	return m, nil
}
