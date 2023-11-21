package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	createAccount(*Account) error
	GetAccountById(int) (*Account, error)
	UpdateAccount(*Account) error
	GetsAllAccount() ([]*Account, error)
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
	fmt.Println("It has connected to database")

	return &PostGresSql{
		db: db,
	}, nil
}
func (p *PostGresSql) init() error {
	return p.createAccoutTable()
}
func (p *PostGresSql) createAccoutTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp,
	)`
	p.db.Exec(query)
	return nil
}
func (p *PostGresSql) createAccount(acct *Account) error {
	fmt.Println(acct)
	query := `insert into account(
		first_name, last_name, number, balance, created_at
	) values($1, $2, $3, $4, $5)`
	_, err := p.db.Query(query, acct.FirstName, acct.LastName, acct.Number, acct.Balance, acct.CreatedAt)
	if err != nil {
		return err
	}
	// fmt.Printf()
	return nil
}
func (p *PostGresSql) GetsAllAccount() ([]*Account, error) {
	query, err := p.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for query.Next() {
		account := new(Account)
		err := query.Scan(&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Balance,
			&account.CreatedAt,
			&account.Number)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
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
