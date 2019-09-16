package jahcoin

import (
	"time"
)

const (
	// MaxTransactionsPerBlock is the amount of transactions per block
	MaxTransactionsPerBlock = 16
)

// Block is a node of a blockchain
type Block struct {
	PrevHash     string
	Hash         string
	Timestamp    time.Time
	Transactions Merkele
}

// Transaction is a transaction between 2 parties
type Transaction struct {
	Amount    float64
	Recipient string
	Sender    string
}

// sign(
// 	hash(Transaction, prevHash),
// 	privateKey
// ) -> SignedTransaction

// az -5
// ti +5

// ti -5
// pesho +5

// sender: na-pesho-det-e-bogat-adresa
// recipeint: moq-adres
// amount: 999999999

// foo(message, publicSenderKey) true/false // verify
// foo(message, privateSenderKey) signedMessage // sign

// foo(message, publicRecipeintKey) message // encrypt
// foo(message, privateRecipeintKey) message // decrypt
