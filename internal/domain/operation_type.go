package domain

// OperationType defines allowed operation/category types.
type OperationType string

const (
	OperationTypeIncome  OperationType = "income"
	OperationTypeExpense OperationType = "expense"
)
