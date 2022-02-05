package viewmodel

// BallanceVM ...
type BallanceVM struct {
	ID          string `json:"id"`
	Amount      int    `json:"amaout"`
	Status      string `json:"status"`
	ReferenceID string `json:"reference_id"`
	DepositedBy string `json:"deposited_by"`
	DepositedAt string `json:"deposited_at"`
	WithdrawnBy string `json:"withdrawn_by"`
	WithdrawnAt string `json:"withdrawn_at"`
}

type WithdrawalVM struct {
	Withdrawal WithdrawalResp `json:"withdrawal"`
}

type WithdrawalResp struct {
	ID          string `json:"id"`
	WithdrawnBy string `json:"withdrawn_by"`
	Status      string `json:"status"`
	WithdrawnAt string `json:"withdrawn_at"`
	Amount      int    `json:"amaout"`
	ReferenceID string `json:"reference_id"`
}

type DepositVM struct {
	Deposit DepositResp `json:"deposit"`
}

type DepositResp struct {
	ID          string `json:"id"`
	DepositedBy string `json:"deposited_by"`
	Status      string `json:"status"`
	DepositedAt string `json:"deposited_at"`
	Amount      int    `json:"amaout"`
	ReferenceID string `json:"reference_id"`
}

type SendQueue struct {
	BalanceID string `json:"balance_id"`
	Amount    int    `json:"amount"`
	OwnedBy   string `json:"owned_by"`
	Type      string `json:"type"`
}
