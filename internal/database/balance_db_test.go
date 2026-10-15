package database

import (
	"database/sql"
	"testing"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type BalanceDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	balanceDB *BalanceDB
}

func (s *BalanceDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	_, err = db.Exec("Create table balances (account_id varchar(255) primary key, balance double, updated_at datetime)")
	s.Nil(err)
	s.balanceDB = NewBalanceDB(db)
}

func (s *BalanceDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE balances")
}

func TestBalanceDBTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceDBTestSuite))
}

func (s *BalanceDBTestSuite) TestSave() {
	balance, _ := entity.NewBalance("account-save", 100.5)
	err := s.balanceDB.Save(balance)
	s.Nil(err)
}

func (s *BalanceDBTestSuite) TestSaveReplacesExistingBalance() {
	first, _ := entity.NewBalance("account-replace", 100)
	s.Nil(s.balanceDB.Save(first))
	second, _ := entity.NewBalance("account-replace", 250.75)
	s.Nil(s.balanceDB.Save(second))

	found, err := s.balanceDB.FindByAccountID("account-replace")
	s.Nil(err)
	s.Equal(250.75, found.Balance)

	var count int
	s.Nil(s.db.QueryRow("SELECT COUNT(*) FROM balances WHERE account_id = ?", "account-replace").Scan(&count))
	s.Equal(1, count)
}

func (s *BalanceDBTestSuite) TestFindByAccountID() {
	balance, _ := entity.NewBalance("account-find", 42.25)
	s.Nil(s.balanceDB.Save(balance))

	found, err := s.balanceDB.FindByAccountID("account-find")
	s.Nil(err)
	s.NotNil(found)
	s.Equal(balance.AccountID, found.AccountID)
	s.Equal(balance.Balance, found.Balance)
	s.False(found.UpdatedAt.IsZero())
}

func (s *BalanceDBTestSuite) TestFindByAccountIDWhenNotFound() {
	found, err := s.balanceDB.FindByAccountID("missing")
	s.Nil(found)
	s.Equal(sql.ErrNoRows, err)
}
