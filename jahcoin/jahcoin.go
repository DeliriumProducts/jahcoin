package jahcoin

import (
	"time"
)

type Block struct {
	PrevHash     string
	Hash         string
	Timestamp    time.Time
	Transactions Merkele
}

type Merkele struct {
	Root         *Node
	transactions []Transaction
}

type Node struct {
	Index    int
	Hash     string
	Children []Node
}

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
