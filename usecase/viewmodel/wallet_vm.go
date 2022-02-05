package viewmodel

type WalletEnableVM struct {
	WalletVM WalletEnableResp `json:"wallet"`
}

// WalletEnableVM ....
type WalletEnableResp struct {
	ID        string `json:"id"`
	OwnedBy   string `json:"owned_by"`
	Status    string `json:"status"`
	EnabledAt string `json:"enabled_at"`
	Balance   int    `json:"balance"`
}

// WalletDisbleVM ...
type WalletDisbleVM struct {
	WalletVM WalletDisbleResp `json:"wallet"`
}

// WalletDisbleResp ....
type WalletDisbleResp struct {
	ID         string `json:"id"`
	OwnedBy    string `json:"owned_by"`
	Status     string `json:"status"`
	DisabledAt string `json:"disabled_at"`
	Balance    int    `json:"balance"`
}

// WalletVM...
type WalletVM struct {
	ID         string `json:"id"`
	Balance    int    `json:"balance"`
	OwnedBy    string `json:"owned_by"`
	Status     string `json:"status"`
	EnabledAt  string `json:"enabled_at"`
	DisabledAt string `json:"disabled_at"`
}
