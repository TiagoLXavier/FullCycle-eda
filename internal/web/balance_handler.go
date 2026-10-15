package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/get_balance"
	"github.com/go-chi/chi/v5"
)

type WebBalanceHandler struct {
	GetBalanceUseCase get_balance.GetBalanceUseCase
}

func NewWebBalanceHandler(getBalanceUseCase get_balance.GetBalanceUseCase) *WebBalanceHandler {
	return &WebBalanceHandler{
		GetBalanceUseCase: getBalanceUseCase,
	}
}

func (h *WebBalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "account_id")
	output, err := h.GetBalanceUseCase.Execute(get_balance.GetBalanceInputDTO{AccountID: accountID})
	if errors.Is(err, sql.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "account not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, output)
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
