package application

import "fmt"

type BankAccountNotFoundError struct {
	Iban string
	Bic  string
}

func NewBankAccountNotFoundError(iban, bic string) error {
	return &BankAccountNotFoundError{iban, bic}
}

func (e *BankAccountNotFoundError) Error() string {
	return fmt.Sprintf("Bank account with IBAN %s AND BIC %s not found", e.Iban, e.Bic)
}

type NotEnoughBalanceError struct {
	Iban string
	Bic  string
}

func NewNotEnoughBalanceError(iban, bic string) error {
	return &NotEnoughBalanceError{iban, bic}
}

func (e *NotEnoughBalanceError) Error() string {
	return fmt.Sprintf("Bank account with IBAN %s AND BIC %s does not have balance enough to perform transfers", e.Iban, e.Bic)
}
