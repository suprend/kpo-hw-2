package model

import "time"

// Payload aggregates exported or imported entities in a transport-friendly form.
type Payload struct {
	Accounts   []Account
	Categories []Category
	Operations []Operation
}

// Account carries bank account data for file operations.
type Account struct {
	ID      string
	Name    string
	Balance int64
}

// Category carries category data for file operations.
type Category struct {
	ID   string
	Type string
	Name string
}

// Operation carries operation data for file operations.
type Operation struct {
	ID            string
	Type          string
	BankAccountID string
	CategoryID    string
	Amount        int64
	Date          time.Time
	Description   string
}
