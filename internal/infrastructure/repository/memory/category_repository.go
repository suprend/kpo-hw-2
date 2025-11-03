package memory

import (
	"sort"
	"strings"
	"sync"

	"kpo-hw-2/internal/domain"
	"kpo-hw-2/internal/domain/repository"
)

type categoryRepository struct {
	mu         sync.RWMutex
	categories map[domain.ID]*domain.Category
}

func NewCategoryRepository() repository.CategoryRepository {
	return &categoryRepository{
		categories: make(map[domain.ID]*domain.Category),
	}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[category.ID()]; exists {
		return domain.ErrAlreadyExists
	}

	r.categories[category.ID()] = category
	return nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[category.ID()]; !exists {
		return domain.ErrNotFound
	}

	r.categories[category.ID()] = category
	return nil
}

func (r *categoryRepository) Delete(id domain.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.categories[id]; !exists {
		return domain.ErrNotFound
	}

	delete(r.categories, id)
	return nil
}

func (r *categoryRepository) Get(id domain.ID) (*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	category, exists := r.categories[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return category, nil
}

func (r *categoryRepository) ListAll() ([]*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.categories) == 0 {
		return nil, nil
	}

	result := make([]*domain.Category, 0, len(r.categories))
	for _, cat := range r.categories {
		result = append(result, cat)
	}

	sortCategories(result)
	return result, nil
}

func (r *categoryRepository) ListByType(typ domain.OperationType) ([]*domain.Category, error) {
	if typ != domain.OperationTypeIncome && typ != domain.OperationTypeExpense {
		return nil, domain.ErrInvalidCategory
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Category
	for _, cat := range r.categories {
		if cat.Type() == typ {
			result = append(result, cat)
		}
	}

	if len(result) == 0 {
		return nil, nil
	}

	sortCategories(result)
	return result, nil
}

func sortCategories(categories []*domain.Category) {
	sort.Slice(categories, func(i, j int) bool {
		ti := categories[i].Type()
		tj := categories[j].Type()
		if ti != tj {
			return ti < tj
		}

		ni := strings.ToLower(categories[i].Name())
		nj := strings.ToLower(categories[j].Name())
		if ni == nj {
			return categories[i].ID() < categories[j].ID()
		}

		return ni < nj
	})
}
