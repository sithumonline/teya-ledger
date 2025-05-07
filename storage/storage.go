package storage

import "sync"

type Storage interface {
	CreateAccount(account Account) (*Account, error)
	GetAccount(userID string, accountNumber string) (*Account, error)

	CreateTransaction(transaction *Transaction) error
	GetTransactions(userID, accountNumber string, limit, page int) ([]*Transaction, error)
	GetTransaction(userID, transactionID string) (*Transaction, error)
	UpdateTransaction(transactionID string, status string) error

	CreateDeposit(transaction *Transaction) (*Transaction, error)
	CreateWithdrawal(transaction *Transaction) (*Transaction, error)

	GetBalance(userID string, accountNumber string) (*Balance, error)
	UpdateBalance(userID string, accountNumber string, amount int64) error
}

type MemoryStorage struct {
	accounts         []*Account
	transactions     []*Transaction
	balances         []*Balance
	accountsLock     sync.RWMutex
	transactionsLock sync.RWMutex
	balancesLock     sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		accounts:     []*Account{},
		transactions: []*Transaction{},
		balances:     []*Balance{},
	}
}
