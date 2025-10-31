package memory

import (
	"sort"
	"strings"
	"sync"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/repository"
)

// accountRepository provides in-memory storage for bank accounts.
type accountRepository struct {
	mu       sync.RWMutex
	accounts map[domain.ID]*domain.BankAccount
}

// NewAccountRepository constructs an in-memory AccountRepository implementation.
func NewAccountRepository() repository.AccountRepository {
	return &accountRepository{
		accounts: make(map[domain.ID]*domain.BankAccount),
	}
}

func (r *accountRepository) Create(account *domain.BankAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.ID()]; exists {
		return domain.ErrAlreadyExists
	}

	r.accounts[account.ID()] = account
	return nil
}

func (r *accountRepository) Update(account *domain.BankAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.ID()]; !exists {
		return domain.ErrNotFound
	}

	r.accounts[account.ID()] = account
	return nil
}

func (r *accountRepository) Delete(id domain.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[id]; !exists {
		return domain.ErrNotFound
	}

	delete(r.accounts, id)
	return nil
}

func (r *accountRepository) Get(id domain.ID) (*domain.BankAccount, error) {
	r.mu.RLock()
	account, exists := r.accounts[id]
	r.mu.RUnlock()

	if !exists {
		return nil, domain.ErrNotFound
	}

	clone := *account
	return &clone, nil
}

func (r *accountRepository) List() ([]*domain.BankAccount, error) {
	r.mu.RLock()
	if len(r.accounts) == 0 {
		r.mu.RUnlock()
		return nil, nil
	}

	result := make([]*domain.BankAccount, 0, len(r.accounts))
	for _, acc := range r.accounts {
		clone := *acc
		result = append(result, &clone)
	}
	r.mu.RUnlock()

	sort.Slice(result, func(i, j int) bool {
		in := strings.ToLower(result[i].Name())
		jn := strings.ToLower(result[j].Name())
		if in == jn {
			return result[i].ID() < result[j].ID()
		}
		return in < jn
	})

	return result, nil
}
