package domain

import (
	"fmt"
)

type TransfersBatch struct {
	OrganizationBic  string     `json:"organization_bic"`
	OrganizationIban string     `json:"organization_iban"`
	OrganizationName string     `json:"organization_name"`
	Transfers        []Transfer `json:"credit_transfers"`
}

type Transfer struct {
	CounterpartyBic  string `json:"counterparty_bic"`
	CounterpartyIban string `json:"counterparty_iban"`
	CounterpartyName string `json:"counterparty_name"`
	// Amount transfer value in EURO
	Amount      float64 `json:"amount,string"`
	Description string  `json:"description,omitempty"`
}

func (b TransfersBatch) Validate() error {
	if len(b.OrganizationBic) == 0 {
		return &EntityValidationError{"organization_bic", "should not be empty"}
	}

	if len(b.OrganizationIban) == 0 {
		return &EntityValidationError{"organization_iban", "should not be empty"}
	}

	if b.Transfers == nil && len(b.Transfers) == 0 {
		return &EntityValidationError{"credit_transfers", "should not be empty"}
	}

	for i, t := range b.Transfers {
		if err := t.validate(); err != nil {
			return &EntityValidationError{fmt.Sprintf("credit_transfers_%d.%s", i, err.Field), err.Message}
		}
	}

	return nil
}

func (b TransfersBatch) Sum() float64 {
	var sum float64 = 0
	for _, t := range b.Transfers {
		sum += t.Amount
	}
	return sum
}

func (b TransfersBatch) SumInCents() int64 {
	var sum int64 = 0
	for _, t := range b.Transfers {
		sum += t.AmountInCents()
	}
	return sum
}

func (t Transfer) validate() *EntityValidationError {
	if t.Amount <= 0 {
		return &EntityValidationError{"amount", "should be more than zero"}
	}

	if len(t.CounterpartyBic) == 0 {
		return &EntityValidationError{"counterparty_bic", "should not be empty"}
	}

	if len(t.CounterpartyIban) == 0 {
		return &EntityValidationError{"counterparty_iban", "should not be empty"}
	}

	return nil
}

func (t Transfer) AmountInCents() int64 {
	return int64(t.Amount * 100)
}
