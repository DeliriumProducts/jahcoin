// jahcoin provides a basic cryptocurrency
package jahcoin

import (
	"errors"
)

var (
	ErrInvalidTransactionsPerBlock = errors.New("jahcoin: TransactionsPerBlock must be a proper logarithm of base 2")
	ErrInvalidDifficulty           = errors.New("jahcoin: Difficulty must be a number higher than 0")
)
