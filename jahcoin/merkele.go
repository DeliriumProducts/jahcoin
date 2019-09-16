package jahcoin

import "errors"

var (
	ErrNotSufficentTransactions = errors.New("merkele: There are not enough transactions to create a merkele tree")
	ErrTooManyTransactions      = errors.New("merkele: There are too many transactions to create a merkele tree")
)

type Node struct {
	Index int
	Hash  string
	Left  *Node
	Right *Node
}

type Merkele struct {
	Root         *Node
	transactions []Transaction
}

// NewMerkele returns a pointer to a Merkele tree and any errors
// The length of the transactions must match the `MaxTransactionsPerBlock` constant
func NewMerkele(transactions []Transaction) (*Merkele, error) {
	tAmount := len(transactions)

	if tAmount >= MaxTransactionsPerBlock {
		return nil, ErrNotSufficentTransactions
	}

	if tAmount < 1 {
		return nil, ErrTooManyTransactions
	}

	m := &Merkele{
		transactions: transactions,
	}

	m.Root = &Node{}

	return m, nil
}
