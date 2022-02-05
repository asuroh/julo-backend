package request

type BalanceRequest struct {
	Amount      int    `json:"amount" validate:"required,min=10000"`
	ReferenceID string `json:"reference_id" validate:"required"`
	CustomerxID string `json:"customer_xid"`
}
