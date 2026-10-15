package database

import (
	"database/sql"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
)

type BalanceDB struct {
	DB *sql.DB
}

func NewBalanceDB(db *sql.DB) *BalanceDB {
	return &BalanceDB{
		DB: db,
	}
}

func (b *BalanceDB) Save(balance *entity.Balance) error {
	stmt, err := b.DB.Prepare("REPLACE INTO balances (account_id, balance, updated_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(balance.AccountID, balance.Balance, balance.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (b *BalanceDB) FindByAccountID(accountID string) (*entity.Balance, error) {
	var balance entity.Balance

	stmt, err := b.DB.Prepare("SELECT account_id, balance, updated_at FROM balances WHERE account_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(accountID)
	err = row.Scan(
		&balance.AccountID,
		&balance.Balance,
		&balance.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}
