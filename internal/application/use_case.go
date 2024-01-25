package application

import (
	"context"
	"transfers-svc/internal/domain"
)

func (a App) CreateTransferBatch(ctx context.Context, batch domain.TransfersBatch) error {
	// For security reasons I am checking the full information provided.
	// So if IBAN is correct but the Name provided does not match, we should not proceed
	accountID, err := a.Repository.GetAccountID(ctx, batch.OrganizationIban, batch.OrganizationBic)
	if err != nil {
		return err
	}

	if accountID == nil || *accountID == 0 {
		return NewBankAccountNotFoundError(batch.OrganizationIban, batch.OrganizationBic)
	}

	// The transaction should be in Read Uncommitted level to allow concurrent processes to have access to the
	// most accurate balance
	transaction, err := a.Repository.BeginReadUncommittedTransaction(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	transferSumIncCents := batch.SumInCents()
	updated, err := a.Repository.DecreaseAccountBalanceOnlyIfBiggerThan(
		ctx,
		transaction,
		*accountID,
		transferSumIncCents,
	)
	if err != nil {
		transaction.Rollback()
		return err
	}

	err = a.Repository.SaveBatch(ctx, transaction, *accountID, batch.Transfers)
	if err != nil {
		transaction.Rollback()
		return err
	}

	if !updated {
		return NewNotEnoughBalanceError(batch.OrganizationIban, batch.OrganizationBic)
	}

	return transaction.Commit()
}
