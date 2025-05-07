package storage

func (m *MemoryStorage) GetBalance(userID string, accountNumber string) (*Balance, error) {
	m.balancesLock.RLock()
	defer m.balancesLock.RUnlock()
	for _, balance := range m.balances {
		if balance.UserID == userID && balance.AccountNumber == accountNumber {
			return balance, nil
		}
	}

	// If no balance exists, create one with 0 amount
	balance := &Balance{
		UserID:        userID,
		AccountNumber: accountNumber,
		Amount:        0,
		Currency:      "MYR",
	}
	m.balances = append(m.balances, balance)
	return balance, nil
}

func (m *MemoryStorage) UpdateBalance(userID string, accountNumber string, amount int64) error {
	m.balancesLock.Lock()
	defer m.balancesLock.Unlock()
	balance, err := m.GetBalance(userID, accountNumber)
	if err != nil {
		return err
	}

	balance.Amount += amount
	return nil
}
