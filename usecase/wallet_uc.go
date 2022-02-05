package usecase

import (
	"database/sql"
	"errors"
	"julo-backend/helper"
	"julo-backend/model"
	"julo-backend/pkg/logruslogger"
	"julo-backend/server/request"
	"julo-backend/usecase/viewmodel"
	"time"
)

// WalletUC ...
type WalletUC struct {
	*ContractUC
	Tx *sql.Tx
}

func (uc WalletUC) Init(req *request.WalletInitRequest) (res viewmodel.JwtVM, err error) {
	const (
		ctx = "WalletUC.Init"
	)

	m := model.NewWalletModel(uc.DB, uc.Tx)
	ok, err := m.WalletExist(req.CustomerxID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "WalletExist", uc.ReqID)
		return res, err
	}

	if !ok {
		_, err = uc.Create(req.CustomerxID)
		if err != nil {
			logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "Create", uc.ReqID)
			return res, err
		}
	}

	payload := map[string]interface{}{
		"customerx_id": req.CustomerxID,
	}
	jwtUc := JwtUC{ContractUC: uc.ContractUC}
	err = jwtUc.GenerateToken(payload, &res)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "jwt", uc.ReqID)
		return res, errors.New(helper.InternalServer)
	}

	return res, err
}

// Enable ...
func (uc WalletUC) Enable(customerxID string) (res viewmodel.WalletEnableVM, err error) {
	const (
		ctx = "WalletUC.Enable"
	)

	m := model.NewWalletModel(uc.DB, uc.Tx)
	status, err := m.FindStatusByOwen(customerxID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "WalletExist", uc.ReqID)
		return res, err
	}

	if status == helper.StatusEnabled && status != "" {
		return res, errors.New(helper.AlreadyEnabled)
	}

	data, err := uc.Update(&request.WalletUpdateRequest{CustomerxID: customerxID, Status: helper.StatusEnabled, EnabledAt: time.Now().Format(time.RFC3339)})
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "Update", uc.ReqID)
		return res, err
	}

	walletData := viewmodel.WalletEnableResp{
		ID:        data.ID,
		OwnedBy:   data.OwnedBy,
		Status:    helper.StatusEnabled,
		Balance:   data.Balance,
		EnabledAt: data.EnabledAt,
	}
	res = viewmodel.WalletEnableVM{WalletVM: walletData}

	return res, err
}

// Disable ...
func (uc WalletUC) Disable(customerxID string) (res viewmodel.WalletDisbleVM, err error) {
	const (
		ctx = "WalletUC.Disable"
	)

	m := model.NewWalletModel(uc.DB, uc.Tx)
	status, err := m.FindStatusByOwen(customerxID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "WalletExist", uc.ReqID)
		return res, err
	}

	if status == helper.StatusDisabled && status != "" {
		return res, errors.New(helper.Disabled)
	}

	data, err := uc.Update(&request.WalletUpdateRequest{CustomerxID: customerxID, Status: helper.StatusDisabled, DisabledAt: time.Now().Format(time.RFC3339)})
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "Update", uc.ReqID)
		return res, err
	}

	walletData := viewmodel.WalletDisbleResp{
		ID:         data.ID,
		OwnedBy:    data.OwnedBy,
		Status:     helper.StatusDisabled,
		Balance:    data.Balance,
		DisabledAt: data.DisabledAt,
	}
	res = viewmodel.WalletDisbleVM{WalletVM: walletData}

	return res, err
}

// Create ...
func (uc WalletUC) Create(customerxID string) (res viewmodel.WalletEnableVM, err error) {
	const (
		ctx = "WalletUC.Create"
	)

	walletData := viewmodel.WalletEnableResp{OwnedBy: customerxID, Balance: 0}
	res = viewmodel.WalletEnableVM{WalletVM: walletData}
	m := model.NewWalletModel(uc.DB, uc.Tx)
	res.WalletVM.ID, err = m.Store(res)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "Store", uc.ReqID)
		return res, err
	}

	return res, err
}

// Update ...
func (uc WalletUC) Update(req *request.WalletUpdateRequest) (res viewmodel.WalletVM, err error) {
	const (
		ctx = "WalletUC.Update"
	)

	res = viewmodel.WalletVM{
		OwnedBy:    req.CustomerxID,
		Status:     req.Status,
		EnabledAt:  req.EnabledAt,
		DisabledAt: req.DisabledAt,
	}

	m := model.NewWalletModel(uc.DB, uc.Tx)
	res.ID, res.Balance, err = m.Update(res)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "Update", uc.ReqID)
		return res, err
	}

	return res, err
}

// GetWallet ...
func (uc WalletUC) GetWallet(customerxID string) (res viewmodel.WalletEnableVM, err error) {
	const (
		ctx = "WalletUC.FindByOwen"
	)

	res, err = uc.FindByOwen(customerxID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "FindByOwen", uc.ReqID)
		return res, err
	}

	if res.WalletVM.Status != helper.StatusEnabled {
		return res, errors.New(helper.Disabled)
	}

	return res, err
}

// FindByOwen ...
func (uc WalletUC) FindByOwen(customerxID string) (res viewmodel.WalletEnableVM, err error) {
	const (
		ctx = "WalletUC.FindByOwen"
	)

	m := model.NewWalletModel(uc.DB, uc.Tx)
	data, err := m.FindByOwen(customerxID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "FindByOwen", uc.ReqID)
		return res, err
	}

	walletData := viewmodel.WalletEnableResp{
		ID:        data.ID,
		OwnedBy:   data.OwnedBy,
		Status:    data.Status.String,
		Balance:   data.Balance,
		EnabledAt: data.EnabledAt.String,
	}
	res = viewmodel.WalletEnableVM{WalletVM: walletData}

	return res, err
}

// CheckBalance ...
func (uc WalletUC) CheckBalance(customerxID string, amount int) (ok bool, err error) {
	const (
		ctx = "WalletUC.CheckBalance"
	)

	m := model.NewWalletModel(uc.DB, uc.Tx)
	balance, err := m.FindBalanceByOwen(customerxID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "FindBalanceByOwen", uc.ReqID)
		return ok, err
	}

	if amount >= balance {
		logruslogger.Log(logruslogger.InfoLevel, "", ctx, "FindBalanceByOwen", uc.ReqID)
		return ok, errors.New(helper.InsufficientBalance)
	}

	return true, err
}

// AddBalance ...
func (uc WalletUC) AddBalance(req viewmodel.SendQueue) (err error) {
	const (
		ctx = "WalletUC.AddBalance"
	)

	m := model.NewWalletModel(uc.DB, uc.Tx)
	if req.Type == helper.TypeWithdrawal {
		_, err = m.MinusBalance(req.OwnedBy, req.Amount)
		if err != nil {
			logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "MinusBalance", uc.ReqID)
			return err
		}
	} else if req.Type == helper.TypeDeposit {
		_, err = m.PlusBalance(req.OwnedBy, req.Amount)
		if err != nil {
			logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "PlusBalance", uc.ReqID)
			return err
		}
	}

	balanceUc := BalanceUC{ContractUC: uc.ContractUC}
	err = balanceUc.UpdateStatus(req.BalanceID, helper.StatusSuccess)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "UpdateStatus", uc.ReqID)
		return err
	}

	return err
}
