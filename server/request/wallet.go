package request

// WalletInitRequest ...
type WalletInitRequest struct {
	CustomerxID string `json:"customer_xid" validate:"required"`
}

// WalletUpdateRequest ...
type WalletUpdateRequest struct {
	CustomerxID string `json:"customer_xid"`
	Status      string `json:"status"`
	EnabledAt   string `json:"enabled_at"`
	DisabledAt  string `json:"disabled_at"`
}
