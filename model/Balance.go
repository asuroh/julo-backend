package model

import (
	"database/sql"
	"julo-backend/helper"
	"julo-backend/usecase/viewmodel"
)

// balanceModel ...
type balanceModel struct {
	DB *sql.DB
	Tx *sql.Tx
}

// IBalance ...
type IBalance interface {
	ReferenceExist(referenceID string) (bool, error)
	StoreWd(body viewmodel.WithdrawalResp) (string, error)
	StoreDe(body viewmodel.DepositResp) (string, error)
	UpdateStatus(id string, status string) error
}

// BalanceEntity ....
type BalanceEntity struct {
	ID          string         `db:"id"`
	Amount      int            `db:"amaout"`
	Status      string         `db:"status"`
	ReferenceID string         `db:"reference_id"`
	DepositedBy sql.NullString `db:"deposited_by"`
	DepositedAt sql.NullString `db:"deposited_at"`
	WithdrawnBy sql.NullString `db:"withdrawn_by"`
	WithdrawnAt sql.NullString `db:"withdrawn_at"`
}

// NewBalanceModel ...
func NewBalanceModel(db *sql.DB, tx *sql.Tx) IBalance {
	return &balanceModel{DB: db, Tx: tx}
}

func (model balanceModel) ReferenceExist(referenceID string) (bool, error) {
	var id string
	sql := `SELECT "reference_id" FROM "balance" WHERE "reference_id" = $1`
	err := model.DB.QueryRow(sql, referenceID).Scan(&id)
	if err != nil {
		if err.Error() == helper.SQLHandlerErrorRowNull {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func (model balanceModel) StoreWd(body viewmodel.WithdrawalResp) (string, error) {
	var id string
	sql := `INSERT INTO "balance" ("amount", "status", "reference_id", "withdrawn_by", "withdrawn_at") VALUES ($1, $2, $3, $4, $5) returning "id"`

	err := model.DB.QueryRow(sql, body.Amount, body.Status, body.ReferenceID, body.WithdrawnBy, body.WithdrawnAt).Scan(&id)

	return id, err
}

func (model balanceModel) StoreDe(body viewmodel.DepositResp) (string, error) {
	var id string
	sql := `INSERT INTO "balance" ("amount", "status", "reference_id", "deposited_by", "deposited_at") VALUES ($1, $2, $3, $4, $5) returning "id"`

	err := model.DB.QueryRow(sql, body.Amount, body.Status, body.ReferenceID, body.DepositedBy, body.DepositedAt).Scan(&id)

	return id, err
}

func (model balanceModel) UpdateStatus(id string, status string) (err error) {
	sql := `UPDATE "balance" SET "status" = $1 WHERE "id" = $2`
	if model.Tx != nil {
		_, err = model.Tx.Exec(sql, status, id)
		return err
	} else {
		_, err = model.DB.Exec(sql, status, id)
	}

	return err
}
