package request

type BalanceRequest struct {
	Amount      int    `json:"amount" validate:"required"`
	ReferenceID string `json:"reference_id" validate:"required"`
	CustomerxID string `json:"customer_xid"`
}
