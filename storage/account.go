package storage

import (
	"errors"
	"time"
)

// CreateTransaction creates a new transaction
func (m *MemoryStorage) CreateAccount(account Account) (*Account, error) {
	m.accountsLock.Lock()
	defer m.accountsLock.Unlock()
	// Ideally should be handled by a unique constraint
	for _, accountData := range m.accounts {
		if account.UserID == account.UserID && accountData.Number == account.Number {
			return nil, errors.New("account already exists")
		}
	}
	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	m.accounts = append(m.accounts, &account)
	return &account, nil
}

func (m *MemoryStorage) GetAccount(userID string, accountNumber string) (*Account, error) {
	m.accountsLock.RLock()
	defer m.accountsLock.RUnlock()
	for _, account := range m.accounts {
		if account.UserID == userID && account.Number == accountNumber {
			return account, nil
		}
	}

	return nil, ErrNotFound
}
