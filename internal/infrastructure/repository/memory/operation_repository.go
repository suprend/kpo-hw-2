package memory

import (
	"sort"
	"sync"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/query"
	"kpo-hw-2/internal/domain/repository"
)

// operationRepository provides in-memory storage for operations.
type operationRepository struct {
	mu         sync.RWMutex
	operations map[domain.ID]*domain.Operation
}

// NewOperationRepository constructs an in-memory OperationRepository implementation.
func NewOperationRepository() repository.OperationRepository {
	return &operationRepository{
		operations: make(map[domain.ID]*domain.Operation),
	}
}

func (r *operationRepository) Create(operation *domain.Operation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.operations[operation.ID()]; exists {
		return domain.ErrAlreadyExists
	}

	r.operations[operation.ID()] = operation
	return nil
}

func (r *operationRepository) Update(operation *domain.Operation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.operations[operation.ID()]; !exists {
		return domain.ErrNotFound
	}

	r.operations[operation.ID()] = operation
	return nil
}

func (r *operationRepository) Delete(id domain.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.operations[id]; !exists {
		return domain.ErrNotFound
	}

	delete(r.operations, id)
	return nil
}

func (r *operationRepository) Get(id domain.ID) (*domain.Operation, error) {
	r.mu.RLock()
	operation, exists := r.operations[id]
	r.mu.RUnlock()

	if !exists {
		return nil, domain.ErrNotFound
	}

	clone := *operation
	return &clone, nil
}

func (r *operationRepository) ListByFilter(filter query.OperationFilter) ([]*domain.Operation, error) {
	r.mu.RLock()

	accountID := filter.AccountID()
	categoryID := filter.CategoryID()
	typ := filter.Type()
	from, to := filter.Period()

	var result []*domain.Operation
	for _, op := range r.operations {
		if accountID != "" && op.BankAccountID() != accountID {
			continue
		}
		if categoryID != "" && op.CategoryID() != categoryID {
			continue
		}
		if typ != "" && op.Type() != typ {
			continue
		}

		opDate := op.Date()
		if from != nil && opDate.Before(*from) {
			continue
		}
		if to != nil && opDate.After(*to) {
			continue
		}

		clone := *op
		result = append(result, &clone)
	}
	r.mu.RUnlock()

	if len(result) == 0 {
		return nil, nil
	}

	sort.Slice(result, func(i, j int) bool {
		di := result[i].Date()
		dj := result[j].Date()
		if di.Equal(dj) {
			return result[i].ID() < result[j].ID()
		}
		return di.Before(dj)
	})

	return result, nil
}
