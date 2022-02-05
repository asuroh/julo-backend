package handler

import (
	"julo-backend/helper"
	"julo-backend/server/request"
	"julo-backend/usecase"
	"net/http"

	validator "gopkg.in/go-playground/validator.v9"
)

// BalanceHandler ...
type BalanceHandler struct {
	Handler
}

// DepositHandler ...
func (h *BalanceHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
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

	if claim["status"].(string) != helper.StatusEnabled {
		SendBadRequest(w, "Disabled")
		return
	}

	req := request.BalanceRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}
	req.CustomerxID = customerxID
	balanceUc := usecase.BalanceUC{ContractUC: h.ContractUC}
	res, err := balanceUc.Deposit(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}

func (h *BalanceHandler) WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
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

	if claim["status"].(string) != helper.StatusEnabled {
		SendBadRequest(w, helper.Disabled)
		return
	}

	req := request.BalanceRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}
	req.CustomerxID = customerxID
	balanceUc := usecase.BalanceUC{ContractUC: h.ContractUC}
	res, err := balanceUc.Withdrawal(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res)
}
