package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TransferRequest struct {
	TransferTo int `json:"transferTo"`
	Amount     int `json:"amonut"`
}
type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type Account struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	EncryptedPassword string    `-`
	Number            int64     `json:"number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleAccount(firstName, lastName, password, email string) (*Account, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		// ID:        rand.Intn(10000),
		Email:             email,
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encrypt),
		Number:            int64(rand.Intn(10000000000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
