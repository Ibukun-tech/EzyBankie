package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	createAccout(*Account) error
	GetAccountById(int) (*Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
}
type PostGresSql struct {
	db *sql.DB
}

func databaseConnection() (*PostGresSql, error) {
	conStr := "user=postgres dbname=postgres password=ezybankgo sslmode=disable"
	db, _ := sql.Open("postgres", conStr)

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &PostGresSql{
		db: db,
	}, nil
}
func (p *PostGresSql) createAccout(*Account) error {
	return nil
}
func (p *PostGresSql) GetAccountById(value int) (*Account, error) {
	return &Account{}, nil
}
func (p *PostGresSql) UpdateAccount(*Account) error {
	return nil
}
func (p *PostGresSql) DeleteAccount(value int) error {
	return nil
}
