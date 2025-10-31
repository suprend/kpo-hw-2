package facade

import (
	"time"

	"kpo-hw-2/internal/domain"
	domainfactory "kpo-hw-2/internal/domain/factory"
	"kpo-hw-2/internal/domain/repository"
)

// operationFacade is the default OperationFacade implementation.
type operationFacade struct {
	factory    domainfactory.OperationFactory
	operations repository.OperationRepository
	accounts   repository.AccountRepository
	categories repository.CategoryRepository
}

// NewOperationFacade wires dependencies for operation workflows.
func NewOperationFacade(
	operationFactory domainfactory.OperationFactory,
	operationRepo repository.OperationRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) OperationFacade {
	return &operationFacade{
		factory:    operationFactory,
		operations: operationRepo,
		accounts:   accountRepo,
		categories: categoryRepo,
	}
}

func (f *operationFacade) CreateOperation(
	typ domain.OperationType,
	accountID domain.ID,
	categoryID domain.ID,
	amount int64,
	date time.Time,
	description string,
) (*domain.Operation, error) {
	context, err := f.buildOperationContext(
		func() (*domain.Operation, error) {
			return f.factory.Create(typ, accountID, categoryID, amount, date, description)
		},
		accountID,
		categoryID,
	)
	if err != nil {
		return nil, err
	}

	if err := f.applyBalance(context.account, context.operation); err != nil {
		return nil, err
	}

	if err := f.operations.Create(context.operation); err != nil {
		_ = f.revertBalanceWithAccount(context.account, context.operation)
		return nil, err
	}

	return context.operation, nil
}

func (f *operationFacade) UpdateOperation(
	id domain.ID,
	typ domain.OperationType,
	accountID domain.ID,
	categoryID domain.ID,
	amount int64,
	date time.Time,
	description string,
) (*domain.Operation, error) {
	existing, err := f.operations.Get(id)
	if err != nil {
		return nil, err
	}

	context, err := f.buildOperationContext(
		func() (*domain.Operation, error) {
			return f.factory.Rebuild(id, typ, accountID, categoryID, amount, date, description)
		},
		accountID,
		categoryID,
	)
	if err != nil {
		return nil, err
	}

	if err := f.updateBalanceForMove(existing, context.operation); err != nil {
		return nil, err
	}

	if err := f.operations.Update(context.operation); err != nil {
		_ = f.updateBalanceForMove(context.operation, existing)
		return nil, err
	}

	return context.operation, nil
}

func (f *operationFacade) DeleteOperation(id domain.ID) error {
	if id == "" {
		return domain.ErrInvalidOperation
	}

	existing, err := f.operations.Get(id)
	if err != nil {
		return err
	}

	if err := f.revertBalance(existing); err != nil {
		return err
	}

	if err := f.operations.Delete(id); err != nil {
		_ = f.applyBalanceByID(existing)
		return err
	}

	return nil
}

type operationContext struct {
	operation *domain.Operation
	account   *domain.BankAccount
	category  *domain.Category
}

func (f *operationFacade) buildOperationContext(
	builder func() (*domain.Operation, error),
	accountID domain.ID,
	categoryID domain.ID,
) (*operationContext, error) {
	op, err := builder()
	if err != nil {
		return nil, err
	}

	if accountID != "" && op.BankAccountID() != accountID {
		return nil, domain.ErrInvalidOperation
	}

	account, err := f.accounts.Get(op.BankAccountID())
	if err != nil {
		return nil, err
	}

	category, err := f.categories.Get(categoryID)
	if err != nil {
		return nil, err
	}

	if category.Type() != op.Type() {
		return nil, domain.ErrInvalidOperation
	}

	return &operationContext{
		operation: op,
		account:   account,
		category:  category,
	}, nil
}

func (f *operationFacade) applyBalance(account *domain.BankAccount, operation *domain.Operation) error {
	if err := account.ApplyOperation(operation); err != nil {
		return err
	}

	if err := f.accounts.Update(account); err != nil {
		_ = account.RevertOperation(operation)
		return err
	}

	return nil
}

func (f *operationFacade) applyBalanceByID(operation *domain.Operation) error {
	account, err := f.accounts.Get(operation.BankAccountID())
	if err != nil {
		return err
	}
	return f.applyBalance(account, operation)
}

func (f *operationFacade) revertBalance(operation *domain.Operation) error {
	account, err := f.accounts.Get(operation.BankAccountID())
	if err != nil {
		return err
	}
	return f.revertBalanceWithAccount(account, operation)
}

func (f *operationFacade) revertBalanceWithAccount(account *domain.BankAccount, operation *domain.Operation) error {
	if err := account.RevertOperation(operation); err != nil {
		return err
	}

	if err := f.accounts.Update(account); err != nil {
		_ = account.ApplyOperation(operation)
		return err
	}

	return nil
}

func (f *operationFacade) updateBalanceForMove(oldOp, newOp *domain.Operation) error {
	if err := f.revertBalance(oldOp); err != nil {
		return err
	}

	if err := f.applyBalanceByID(newOp); err != nil {
		_ = f.applyBalanceByID(oldOp)
		return err
	}

	return nil
}

func (f *operationFacade) ListOperations(
	accountID domain.ID,
	from time.Time,
	to time.Time,
) ([]*domain.Operation, error) {
	if accountID == "" {
		return nil, domain.ErrInvalidOperation
	}

	if from.After(to) {
		return nil, domain.ErrInvalidOperation
	}

	return f.operations.ListByAccountAndPeriod(accountID, from, to)
}

func (f *operationFacade) GetOperation(id domain.ID) (*domain.Operation, error) {
	if id == "" {
		return nil, domain.ErrInvalidOperation
	}

	return f.operations.Get(id)
}

var _ OperationFacade = (*operationFacade)(nil)
