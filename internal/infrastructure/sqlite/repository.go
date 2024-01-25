package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"transfers-svc/internal/domain"

	"github.com/gocraft/dbr/v2"
)

const bankAccountsTable = "bank_accounts"
const transfersTable = "transfers"

type Repository struct {
	DbSession *dbr.Session
}

func (r Repository) SaveBatch(
	ctx context.Context,
	transaction domain.Transaction,
	bankAccountID int64,
	transfers []domain.Transfer,
) error {
	tx, ok := transaction.(Transaction)
	if !ok {
		return errors.New("invalid transaction")
	}

	stmt := tx.tx.InsertInto(transfersTable).Columns(
		"bank_account_id",
		"counterparty_bic",
		"counterparty_iban",
		"counterparty_name",
		"amount_cents",
		"description")

	for _, t := range transfers {
		row := TransferRow{
			BankAccountID:    bankAccountID,
			CounterpartyBic:  t.CounterpartyBic,
			CounterpartyIban: t.CounterpartyName,
			CounterpartyName: t.CounterpartyIban,
			AmountCents:      t.AmountInCents(),
			Description:      t.Description,
		}
		stmt.Record(row)
	}

	result, err := stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(transfers)) {
		return errors.New("not all transfers were persisted in DB")
	}

	return nil
}

func (r Repository) BeginReadUncommittedTransaction(ctx context.Context) (domain.Transaction, error) {
	tx, err := r.DbSession.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{tx}, err
}

func (r Repository) GetAccountID(ctx context.Context, iban, bic string) (*int64, error) {
	var accountID int64

	err := r.DbSession.
		Select("id").
		From(bankAccountsTable).
		Where(dbr.Expr("iban = ? AND bic = ?", iban, bic)).
		LoadOneContext(ctx, &accountID)

	if err == dbr.ErrNotFound {
		return nil, nil
	}

	return &accountID, err
}

func (r Repository) DecreaseAccountBalanceOnlyIfBiggerThan(
	ctx context.Context,
	transaction domain.Transaction,
	bankAccountID int64,
	transfersTotalInCents int64,
) (bool, error) {
	tx, ok := transaction.(Transaction)
	if !ok {
		return false, errors.New("invalid transaction")
	}

	query := fmt.Sprintf(
		"UPDATE %s SET balance_cents = balance_cents - %d WHERE id = %d AND balance_cents >= %d",
		bankAccountsTable,
		transfersTotalInCents,
		bankAccountID,
		transfersTotalInCents)

	result, err := tx.tx.ExecContext(ctx, query)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
