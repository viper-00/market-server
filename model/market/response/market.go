package response

import "market/model"

type NetworkInfo struct {
	TatumUrl string `json:"tatum_url"`
	TatumKey string `json:"tatum_key"`
	ChainId  int    `json:"chain_id"`
	RPCUrl   string `json:"rpc_url"`
	HTTPUrl  string `json:"http_url"`
	HTTPKey  string `json:"http_key"`
}

type StoreUserWallet struct {
	Address string `json:"address" `
	ChainId int    `json:"chain_id"`
}

type BulkStoreUserWalletResponse struct {
	BulkStorage []StoreUserWallet `json:"bulk_storage"`
}

type StoreChainContract struct {
	Contract string `json:"contract" `
	ChainId  int    `json:"chain_id"`
}

type BulkStoreChainContractResponse struct {
	BulkStorage []StoreChainContract `json:"bulk_storage"`
}

type OwnListResponse struct {
	Transactions []model.OwnTransaction `json:"transactions"`
	Total        int64                  `json:"total"`
	Page         int                    `json:"page"`
	PageSize     int                    `json:"pageSize"`
}

type EventPlayResponse struct {
	Title              string                   `json:"title"`
	Introduce          string                   `json:"introduce"`
	GuessNumber        int                      `json:"guess_number"`
	MinimumCapitalPool float64                  `json:"minimum_capital_pool"`
	MaximumCapitalPool float64                  `json:"maximum_capital_pool"`
	Coin               string                   `json:"coin"`
	PledgeAmount       float64                  `json:"pledge_amount"`
	Values             []EventPlayValueResponse `json:"values"`
}

type EventPlayValueResponse struct {
	Value  string               `json:"value"`
	Orders []EventOrderResponse `json:"orders"`
}

type EventOrderResponse struct {
	Amount              float64 `json:"amount"`
	OrderType           string  `json:"order_type"`
	UserContractAddress string  `json:"user_address"`
	Username            string  `json:"username"`
	CreatedTime         int     `json:"created_time"`
	Hash                string  `json:"hash"`
}

type EventCommentResponse struct {
	Username            string `json:"username"`
	AvatarUrl           string `json:"avatar_url"`
	Content             string `json:"content"`
	CommentId           uint   `json:"comment_id"`
	CreatedTime         int    `json:"created_time"`
	UserContractAddress string `json:"user_address"`
	LikeCount           int    `json:"like_count"`
	OwnLikeStatus       int    `json:"own_like_status"`
}

type EventForHomeResponse struct {
	EventLogo        string  `json:"event_logo"`
	Title            string  `json:"title"`
	ExpireTime       int     `json:"expire_time"`
	UniqueCode       string  `json:"unique_code"`
	Type             string  `json:"type"`
	SettlementTime   int     `json:"settlement_time"`
	TotalOrderAmount float64 `json:"total_order_amount"`
	CommentCount     int     `json:"comment_count"`
}

type TopVolumnForHomeResponse struct {
	AvatarUrl           string `json:"avatar_url"`
	Username            string `json:"username"`
	CryptoAmount        string `json:"crypto_amount"`
	LegalAmount         string `json:"legal_amount"`
	UserContractAddress string `json:"user_address"`
}

type RecentActivityForHomeResponse struct {
	EventLogo   string  `json:"event_logo"`
	Title       string  `json:"title"`
	UniqueCode  string  `json:"unique_code"`
	CreatedTime int     `json:"created_time"`
	AvatarUrl   string  `json:"avatar_url"`
	Amount      float64 `json:"amount"`
	OrderType   string  `json:"order_type"`
	Username    string  `json:"username"`
	PlayValue   string  `json:"play_value"`
}
