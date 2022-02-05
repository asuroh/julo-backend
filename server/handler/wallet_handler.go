package handler

import (
	"julo-backend/helper"
	"julo-backend/server/request"
	"julo-backend/usecase"
	"net/http"

	validator "gopkg.in/go-playground/validator.v9"
)

// WalletHandler ...
type WalletHandler struct {
	Handler
}

// InitHandler ...
func (h *WalletHandler) InitHandler(w http.ResponseWriter, r *http.Request) {
	req := request.WalletInitRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}
	walletUc := usecase.WalletUC{ContractUC: h.ContractUC}
	res, err := walletUc.Init(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}

// EnableHandler ...
func (h *WalletHandler) EnableHandler(w http.ResponseWriter, r *http.Request) {
	claim := requestIDFromContextInterface(r.Context(), helper.Token)
	if claim == nil {
		SendBadRequest(w, "Invalid claim")
		return
	}

	customerxID := claim["customerx_id"].(string)
	if customerxID == "" {
		SendBadRequest(w, "Invalid customerx id")
		return
	}

	walletUc := usecase.WalletUC{ContractUC: h.ContractUC}
	res, err := walletUc.Enable(customerxID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}

// GetWalletHandler ...
func (h *WalletHandler) GetWalletHandler(w http.ResponseWriter, r *http.Request) {
	claim := requestIDFromContextInterface(r.Context(), helper.Token)
	if claim == nil {
		SendBadRequest(w, "Invalid claim")
		return
	}

	customerxID := claim["customerx_id"].(string)
	if customerxID == "" {
		SendBadRequest(w, "Invalid customerx id")
		return
	}

	walletUc := usecase.WalletUC{ContractUC: h.ContractUC}
	res, err := walletUc.GetWallet(customerxID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}

// DisableHandler ...
func (h *WalletHandler) DisableHandler(w http.ResponseWriter, r *http.Request) {
	claim := requestIDFromContextInterface(r.Context(), helper.Token)
	if claim == nil {
		SendBadRequest(w, "Invalid claim")
		return
	}

	customerxID := claim["customerx_id"].(string)
	if customerxID == "" {
		SendBadRequest(w, "Invalid customerx id")
		return
	}

	walletUc := usecase.WalletUC{ContractUC: h.ContractUC}
	res, err := walletUc.Disable(customerxID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}
