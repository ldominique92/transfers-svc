package application_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"transfers-svc/internal/infrastructure/sqlite"

	"transfers-svc/internal/application"
	"transfers-svc/internal/domain"

	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

const (
	organizationIBAN = "NL12345678901234567890XXXXX"
	organizationBIC  = "ABCDEFGHIJK"
)

func TestTransfersHandler(t *testing.T) {
	// Configure DB connection
	dbConn, err := dbr.Open("sqlite3", "../../qonto_accounts.sqlite", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer dbConn.Close()

	dbSession := dbConn.NewSession(nil)
	if err != nil {
		t.Fatal(err)
	}

	testApp := &application.App{
		Repository: sqlite.Repository{DbSession: dbSession},
	}

	// Reset DB
	_, err = dbSession.DeleteFrom("bank_accounts").Where("iban = ?", organizationIBAN).Exec()
	if err != nil {
		t.Fatal(err)
	}

	// With invalid transfer batch, should return 400
	transferBatch := &domain.TransfersBatch{}
	rr := executeRequest(t, transferBatch, testApp)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, strings.Trim(rr.Body.String(), "\n"), "organization_bic : should not be empty")

	// With valid transfer batch
	transferBatch = &domain.TransfersBatch{
		OrganizationBic:  organizationBIC,
		OrganizationIban: organizationIBAN,
		OrganizationName: "Test Org",
		Transfers: []domain.Transfer{
			{
				CounterpartyBic:  "LMNOPRSTUVW",
				CounterpartyIban: "FR12345678901234567890XXXXX",
				CounterpartyName: "Another Org",
				Amount:           100.52,
				Description:      "Lorem ipsum",
			},
			{
				CounterpartyBic:  "XYZABCDEFGH",
				CounterpartyIban: "DE12345678901234567890XXXXX",
				CounterpartyName: "One more Org",
				Amount:           53.08,
				Description:      "Lorem ipsum",
			},
		},
	}

	// When account does not exist, it should return HTTP status 404
	rr = executeRequest(t, transferBatch, testApp)
	assert.Equal(t, rr.Code, http.StatusNotFound)
	assert.Equal(t, strings.Trim(rr.Body.String(), "\n"), fmt.Sprintf(
		"Bank account with IBAN %s AND BIC %s not found",
		organizationIBAN,
		organizationBIC))

	// When user account exists and the user has funds, it should return HTTP status 201
	_, err = dbSession.InsertInto("bank_accounts").
		Pair("organization_name", "Test Org").
		Pair("balance_cents", 20000).
		Pair("iban", organizationIBAN).
		Pair("bic", organizationBIC).
		Exec()

	if err != nil {
		t.Fatal(err)
	}

	rr = executeRequest(t, transferBatch, testApp)
	assert.Equal(t, rr.Code, http.StatusCreated)

	// And the account balance should be updated correctly
	var accountBalance int
	err = dbSession.Select("balance_cents").
		From("bank_accounts").
		Where("iban = ?", organizationIBAN).
		LoadOne(&accountBalance)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 4640, accountBalance)

	// TODO: we should assert that the transfers were created correctly

	// When account does not have balance enough, should return 422
	rr = executeRequest(t, transferBatch, testApp)
	assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
	assert.Equal(t, strings.Trim(rr.Body.String(), "\n"),
		fmt.Sprintf(
			"Bank account with IBAN %s AND BIC %s does not have balance enough to perform transfers",
			organizationIBAN,
			organizationBIC))
}

func executeRequest(t *testing.T, transferBatch *domain.TransfersBatch, mockApp *application.App) *httptest.ResponseRecorder {
	jsonData, _ := json.Marshal(transferBatch)
	req, err := http.NewRequest("POST", "/transfers", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mockApp.TransfersHandler)
	handler.ServeHTTP(rr, req)
	return rr
}
