package sqlite

type TransferRow struct {
	ID               int64  `db:"id"`
	BankAccountID    int64  `db:"bank_account_id"`
	CounterpartyBic  string `db:"counterparty_bic"`
	CounterpartyIban string `db:"counterparty_iban"`
	CounterpartyName string `db:"counterparty_name"`
	AmountCents      int64  `db:"amount_cents"`
	Description      string `db:"description"`
}
