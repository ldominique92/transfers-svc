package application

import (
	"encoding/json"
	"net/http"
	"transfers-svc/internal/domain"

	"github.com/labstack/gommon/log"
)

func (a App) TransfersHandler(w http.ResponseWriter, r *http.Request) {
	var transferBatch domain.TransfersBatch
	err := json.NewDecoder(r.Body).Decode(&transferBatch)
	if err != nil {
		log.Error(err)
		http.Error(w, "Request was not in expected format", http.StatusBadRequest)
		return
	}

	if err := transferBatch.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.CreateTransferBatch(r.Context(), transferBatch); err != nil {
		switch t := err.(type) {
		case *BankAccountNotFoundError:
			http.Error(w, t.Error(), http.StatusNotFound)
			return
		case *NotEnoughBalanceError:
			http.Error(w, t.Error(), http.StatusUnprocessableEntity)
			return
		}

		log.Errorf("Unexpected error while trying transfer money from account %s: %s", transferBatch.OrganizationIban, err)
		http.Error(w, "Unexpected error while trying transfer money from account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
