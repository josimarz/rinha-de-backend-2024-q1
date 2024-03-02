package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/josimarz/rinha-de-backend-2024-q1/internal/entity"
	"github.com/josimarz/rinha-de-backend-2024-q1/internal/usecase"
)

type BankStatementHandler struct {
	uc *usecase.GetBankStatementUseCase
}

func NewBankStatementHandler(uc *usecase.GetBankStatementUseCase) *BankStatementHandler {
	return &BankStatementHandler{uc}
}

func (h *BankStatementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := r.PathValue("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	o, err := h.uc.Execute(uint8(id))
	if err != nil {
		if _, ok := err.(*usecase.AccountNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(o); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type TransactionBody struct {
	Amount int64  `json:"valor"`
	Kind   string `json:"tipo"`
	Desc   string `json:"descricao"`
}

type TransactionHandler struct {
	uc *usecase.DoTransactionUseCase
}

func NewTransactionHandler(uc *usecase.DoTransactionUseCase) *TransactionHandler {
	return &TransactionHandler{uc}
}

func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := r.PathValue("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := h.parseBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	i := &usecase.DoTransactionInput{
		AccountId: uint8(id),
		Amount:    body.Amount,
		Kind:      body.Kind,
		Desc:      body.Desc,
	}
	o, err := h.uc.Execute(i)
	if err != nil {
		if _, ok := err.(*usecase.AccountNotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if _, ok := err.(*entity.InsuficientBalanceError); ok {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(o); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TransactionHandler) parseBody(r io.Reader) (*TransactionBody, error) {
	body := &TransactionBody{}
	if err := json.NewDecoder(r).Decode(body); err != nil {
		return nil, err
	}
	if body.Amount <= 0 {
		return nil, errors.New("amount should be a positive value")
	}
	if body.Kind != "c" && body.Kind != "d" {
		return nil, errors.New("invalid operation")
	}
	if l := len(body.Desc); l < 1 || l > 10 {
		return nil, errors.New("description out of bounds")
	}
	return body, nil
}
