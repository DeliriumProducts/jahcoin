package main

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/deliriumproducts/jahcoin/jahcoin"
)

func main() {
	gekPub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	bc, err := jahcoin.NewBlockchain(&jahcoin.Config{
		TransactionsPerBlock:     16,
		Difficulty:               5,
		Gekyume:                  gekPub,
		InitialTransactionAmount: 50,
	})

	if err != nil {
		panic(err)
	}

	bc.AddTransaction(&jahcoin.Transaction{
		Amount: 999,
	})
}
