package model

import "time"

type Payload struct {
	Accounts   []Account
	Categories []Category
	Operations []Operation
}

type Account struct {
	ID      string
	Name    string
	Balance int64
}

type Category struct {
	ID   string
	Type string
	Name string
}

type Operation struct {
	ID            string
	Type          string
	BankAccountID string
	CategoryID    string
	Amount        int64
	Date          time.Time
	Description   string
}
