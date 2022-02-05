package usecase

import (
	"database/sql"
	"errors"
	"julo-backend/helper"
	"julo-backend/model"
	"julo-backend/pkg/amqp"
	"julo-backend/pkg/logruslogger"
	"julo-backend/server/request"
	"julo-backend/usecase/viewmodel"
	"time"
)

// BalanceUC ...
type BalanceUC struct {
	*ContractUC
	Tx *sql.Tx
}

func (uc *BalanceUC) Deposit(req *request.BalanceRequest) (res viewmodel.DepositVM, err error) {
	const (
		ctx = "Deposit"
	)

	m := model.NewBalanceModel(uc.DB, uc.Tx)
	ok, err := m.ReferenceExist(req.ReferenceID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "ReferenceExist", uc.ReqID)
		return res, err
	}

	if ok {
		logruslogger.Log(logruslogger.InfoLevel, "", ctx, "ReferenceExist", uc.ReqID)
		return res, errors.New(helper.ReferenceExist)
	}

	now := time.Now().Format(time.RFC3339)
	res.Deposit = viewmodel.DepositResp{
		Amount:      req.Amount,
		Status:      helper.StatusPending,
		ReferenceID: req.ReferenceID,
		DepositedBy: req.CustomerxID,
		DepositedAt: now,
	}

	res.Deposit.ID, err = m.StoreDe(res.Deposit)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "StoreDe", uc.ReqID)
		return res, err
	}

	err = uc.sendQueue(viewmodel.SendQueue{
		OwnedBy:   req.CustomerxID,
		Amount:    req.Amount,
		Type:      helper.TypeDeposit,
		BalanceID: res.Deposit.ID,
	})
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "sendQueue", uc.ReqID)
		return res, err
	}

	return res, err
}

func (uc *BalanceUC) Withdrawal(req *request.BalanceRequest) (res viewmodel.WithdrawalVM, err error) {
	const (
		ctx = "Withdrawal"
	)

	m := model.NewBalanceModel(uc.DB, uc.Tx)
	ok, err := m.ReferenceExist(req.ReferenceID)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "ReferenceExist", uc.ReqID)
		return res, err
	}

	if ok {
		logruslogger.Log(logruslogger.InfoLevel, "", ctx, "ReferenceExist", uc.ReqID)
		return res, errors.New(helper.ReferenceExist)
	}

	walletUc := WalletUC{ContractUC: uc.ContractUC}
	ok, err = walletUc.CheckBalance(req.CustomerxID, req.Amount)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "FindBalance", uc.ReqID)
		return res, err
	}

	if !ok {
		logruslogger.Log(logruslogger.InfoLevel, "", ctx, "FindBalance", uc.ReqID)
		return res, err
	}

	now := time.Now().Format(time.RFC3339)
	res.Withdrawal = viewmodel.WithdrawalResp{
		Amount:      req.Amount,
		Status:      helper.StatusPending,
		ReferenceID: req.ReferenceID,
		WithdrawnBy: req.CustomerxID,
		WithdrawnAt: now,
	}

	res.Withdrawal.ID, err = m.StoreWd(res.Withdrawal)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "StoreWd", uc.ReqID)
		return res, err
	}

	err = uc.sendQueue(viewmodel.SendQueue{
		OwnedBy:   req.CustomerxID,
		Amount:    req.Amount,
		Type:      helper.TypeWithdrawal,
		BalanceID: res.Withdrawal.ID,
	})
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "sendQueue", uc.ReqID)
		return res, err
	}

	return res, err
}

func (uc BalanceUC) UpdateStatus(id, status string) (err error) {
	const (
		ctx = "UpdateStatus"
	)

	m := model.NewBalanceModel(uc.DB, uc.Tx)
	err = m.UpdateStatus(id, status)
	if err != nil {
		logruslogger.Log(logruslogger.ErrorLevel, err.Error(), ctx, "UpdateStatus", uc.ReqID)
		return err
	}

	return err
}

func (uc BalanceUC) sendQueue(req viewmodel.SendQueue) (err error) {
	const (
		ctx = "sendQueue"
	)

	mqueue := amqp.NewQueue(AmqpConnection, AmqpChannel)
	queueBody := map[string]interface{}{
		"qid":        uc.ContractUC.ReqID,
		"owned_by":   req.OwnedBy,
		"amount":     req.Amount,
		"type":       req.Type,
		"balance_id": req.BalanceID,
	}
	AmqpConnection, AmqpChannel, err = mqueue.PushQueueReconnect(uc.ContractUC.EnvConfig["AMQP_URL"], queueBody, amqp.UpdateBalance, amqp.UpdateBalanceDeadLetter)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "update_balance_queue", uc.ReqID)
		return errors.New("update_balance_queue")
	}

	return err
}
