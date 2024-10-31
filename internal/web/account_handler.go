package web

import (
	"encoding/json"
	"net/http"

	"github.com/wallison/fc-ms-wallet/internal/usecase/create_account"
)

type WebAccontHandler struct {
	CreateAccountUseCase create_account.CreateAccountUseCase
}

func NewWebAccountHandler(createAccountUseCase create_account.CreateAccountUseCase) *WebAccontHandler {
	return &WebAccontHandler{
		CreateAccountUseCase: createAccountUseCase,
	}
}

func (h WebAccontHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var dto create_account.CreateAccountInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.CreateAccountUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
