package main

import "database/sql"

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
	return nil, nil
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
