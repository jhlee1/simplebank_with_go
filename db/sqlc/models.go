// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package simplebank

import (
	"database/sql"
	"time"
)

type Account struct {
	ID        int64
	Owner     string
	Balance   int64
	Currency  string
	CreatedAt time.Time
}

type Entry struct {
	ID        int64
	AccountID sql.NullInt64
	// can be negative or positive
	Amount    int64
	CreatedAt time.Time
}

type Transfer struct {
	ID            int64
	FromAccountID sql.NullInt64
	ToAccountID   sql.NullInt64
	// must be positive
	Amount    int64
	CreatedAt sql.NullTime
}
