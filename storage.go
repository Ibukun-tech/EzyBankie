package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	FindAccount(*TransferRequest) (*TransferRequest, error)
	GetAccountByNumber(int) (*Account, error)
	createAccount(*Account) error
	GetAccountById(int) (*Account, error)
	UpdateAccount(*Account) error
	GetsAllAccount() ([]*Account, error)
	DeleteAccount(int) error
}
type PostGresSql struct {
	db *sql.DB
}

func (p *PostGresSql) FindAccount(tr *TransferRequest) (*TransferRequest, error) {
	query, err := p.db.Query("select * from accounts where number =$1", tr.Number)
	if err != nil {
		return nil, err
	}
	transfer := &TransferRequest{}
	err = query.Scan(&transfer.Number)
	if err != nil {
		return nil, err
	}
	if transfer.Number != tr.Number {
		return nil, fmt.Errorf("its not the same account number")
	}
	return nil, nil
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
		email varchar(60)
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp,
	)`
	p.db.Exec(query)
	return nil
}
func (p *PostGresSql) createAccount(acct *Account) error {
	// fmt.Println(acct)
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	query := `insert into account(
		first_name, last_name, number, balance, created_at
	) values($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, acct.FirstName, acct.LastName, acct.Number, acct.Balance, acct.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return nil
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
		account, err := scanIntoAccount(query)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID, &account.FirstName,
		&account.LastName,
		&account.Balance,
		&account.CreatedAt,
		&account.Number)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (p *PostGresSql) GetAccountByNumber(number int) (*Account, error) {
	rows, err := p.db.Query("select * from account where number = $1", number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", number)
}
func (p *PostGresSql) GetAccountById(value int) (*Account, error) {
	rows, err := p.db.Query("select * from account where id = $1", value)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", value)
}
func (p *PostGresSql) UpdateAccount(*Account) error {
	return nil
}
func (p *PostGresSql) DeleteAccount(value int) error {
	_, err := p.db.Query("delete from account where id = $1", value)
	if err != nil {
		return err
	}
	return nil
}
