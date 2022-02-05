package model

import (
	"database/sql"
	"fmt"
	"julo-backend/helper"
	"julo-backend/usecase/viewmodel"
)

// walletModel ...
type walletModel struct {
	DB *sql.DB
	Tx *sql.Tx
}

// IWallet ...
type IWallet interface {
	WalletExist(ownedBy string) (bool, error)
	FindStatusByOwen(ownedBy string) (string, error)
	FindBalanceByOwen(ownedBy string) (int, error)
	FindByOwen(ownedBy string) (WalletEntity, error)
	Store(body viewmodel.WalletEnableVM) (string, error)
	Update(body viewmodel.WalletVM) (string, int, error)
	MinusBalance(ownedBy string, amount int) (string, error)
	PlusBalance(ownedBy string, amount int) (string, error)
}

// WalletEntity ....
type WalletEntity struct {
	ID         string         `db:"id"`
	Balance    int            `db:"balance"`
	OwnedBy    string         `db:"owned_by"`
	Status     sql.NullString `db:"status"`
	EnabledAt  sql.NullString `db:"enabled_at"`
	DisabledAt sql.NullString `db:"disabled_at"`
}

// NewWalletModel ...
func NewWalletModel(db *sql.DB, tx *sql.Tx) IWallet {
	return &walletModel{DB: db, Tx: tx}
}

func (model walletModel) WalletExist(ownedBy string) (bool, error) {
	var id sql.NullString
	sql := `SELECT "id" FROM "wallet" WHERE "owned_by" = $1`
	err := model.DB.QueryRow(sql, ownedBy).Scan(&id)
	if err != nil {
		if err.Error() == helper.SQLHandlerErrorRowNull {
			return false, nil
		}

		return false, err
	}

	return true, err
}

func (model walletModel) FindStatusByOwen(ownedBy string) (string, error) {
	var status sql.NullString
	sql := `SELECT "status" FROM "wallet" WHERE "owned_by" = $1`
	err := model.DB.QueryRow(sql, ownedBy).Scan(&status)
	if err != nil {
		if err.Error() == helper.SQLHandlerErrorRowNull {
			return "", nil
		}

		return "", err
	}

	return status.String, err
}

func (model walletModel) FindBalanceByOwen(ownedBy string) (int, error) {
	var balance int
	sql := `SELECT "balance" FROM "wallet" WHERE "owned_by" = $1`
	err := model.DB.QueryRow(sql, ownedBy).Scan(&balance)
	if err != nil {
		if err.Error() == helper.SQLHandlerErrorRowNull {
			return 0, nil
		}

		return 0, err
	}

	return balance, err
}

func (model walletModel) FindByOwen(ownedBy string) (WalletEntity, error) {
	var d WalletEntity
	sql := `SELECT "id", "balance", "owned_by", "status", "enabled_at", "disabled_at" FROM "wallet" WHERE "owned_by" = $1`
	err := model.DB.QueryRow(sql, ownedBy).Scan(
		&d.ID, &d.Balance, &d.OwnedBy, &d.Status,
		&d.EnabledAt, &d.DisabledAt,
	)
	if err != nil {
		if err.Error() == helper.SQLHandlerErrorRowNull {
			return d, nil
		}

		return d, err
	}

	return d, err
}

// Store ...
func (model walletModel) Store(body viewmodel.WalletEnableVM) (res string, err error) {
	sql := `INSERT INTO "wallet" (
			"balance", "owned_by"
		) VALUES($1, $2) RETURNING "id"`
	err = model.DB.QueryRow(sql, body.WalletVM.Balance, body.WalletVM.OwnedBy).Scan(&res)

	return res, err
}

// Update ...
func (model walletModel) Update(body viewmodel.WalletVM) (id string, balance int, err error) {
	sql := `UPDATE "wallet" SET "status" = $1, "disabled_at" = $2, "enabled_at" = $3 WHERE "owned_by" = $4 RETURNING "id", "balance"`
	err = model.DB.QueryRow(sql, body.Status, newNullString(body.DisabledAt), newNullString(body.EnabledAt), body.OwnedBy).Scan(&id, &balance)

	return id, balance, err
}

// MinusBalance ...
func (model walletModel) MinusBalance(ownedBy string, amount int) (res string, err error) {
	sql := `UPDATE "wallet" SET "balance" = "balance" - $1 WHERE "owned_by" = $2 RETURNING "id"`
	if model.Tx != nil {
		err = model.Tx.QueryRow(sql, amount, ownedBy).Scan(&res)
	} else {
		err = model.DB.QueryRow(sql, amount, ownedBy).Scan(&res)
	}

	return res, err
}

// PlusBalance ...
func (model walletModel) PlusBalance(ownedBy string, amount int) (res string, err error) {
	sql := `UPDATE "wallet" SET "balance" = "balance" + $1 WHERE "owned_by" = $2 RETURNING "id"`
	if model.Tx != nil {
		fmt.Println("sampe sini")
		err = model.Tx.QueryRow(sql, amount, ownedBy).Scan(&res)
	} else {
		err = model.DB.QueryRow(sql, amount, ownedBy).Scan(&res)
	}

	return res, err
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
