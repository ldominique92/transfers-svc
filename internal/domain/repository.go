package domain

import "context"

type Repository interface {
	SaveBatch(ctx context.Context, transaction Transaction, bankAccountID int64, transfers []Transfer) error
	BeginReadUncommittedTransaction(ctx context.Context) (Transaction, error)
	GetAccountID(ctx context.Context, iban, bic string) (*int64, error)
	DecreaseAccountBalanceOnlyIfBiggerThan(
		ctx context.Context,
		transaction Transaction,
		bankAccountID int64,
		transfersSumInCents int64,
	) (bool, error)
}

type Transaction interface {
	Rollback()
	Commit() error
}
