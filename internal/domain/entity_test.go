package domain_test

import (
	"testing"
	"transfers-svc/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestTransfer_AmountInCents(t *testing.T) {
	transfer := domain.Transfer{
		CounterpartyBic:  "test",
		CounterpartyIban: "test",
		CounterpartyName: "test",
		Description:      "test",
	}

	transfer.Amount = 55
	assert.Equal(t, transfer.AmountInCents(), int64(5500))

	transfer.Amount = 58.95
	assert.Equal(t, transfer.AmountInCents(), int64(5895))

	transfer.Amount = 65.7
	assert.Equal(t, transfer.AmountInCents(), int64(6570))

	transfer.Amount = 0.29
	assert.Equal(t, transfer.AmountInCents(), int64(29))
}
