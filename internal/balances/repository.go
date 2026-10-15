package balances

import (
	"database/sql"
	"time"
)

type Balance struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Upsert(accountID string, balance float64) error {
	_, err := r.DB.Exec(
		"INSERT INTO balances (account_id, balance, updated_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE balance = VALUES(balance), updated_at = VALUES(updated_at)",
		accountID, balance, time.Now(),
	)
	return err
}

func (r *Repository) FindByAccountID(accountID string) (*Balance, error) {
	var b Balance
	err := r.DB.QueryRow("SELECT account_id, balance FROM balances WHERE account_id = ?", accountID).Scan(&b.AccountID, &b.Balance)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
