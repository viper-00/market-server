package request

type StoreUserWallet struct {
	Address string `json:"address" form:"address" binding:"required"`
	ChainId int    `json:"chain_id" form:"chain_id" binding:"required"`
}

type BulkStoreUserWallet struct {
	BulkStorage []StoreUserWallet `json:"bulk_storage" form:"bulk_storage" binding:"required"`
}

type GetNetworkInfo struct {
	ChainId int `json:"chain_id" form:"chain_id" binding:"required"`
}

type StoreChainContract struct {
	ChainId  int    `json:"chain_id" form:"chain_id" binding:"required"`
	Symbol   string `json:"symbol" form:"symbol" binding:"required"`
	Decimals int    `json:"decimals" form:"decimals" binding:"required"`
	Contract string `json:"contract" form:"contract" binding:"required"`
}

type BulkStoreChainContract struct {
	BulkStorage []StoreChainContract `json:"bulk_storage" form:"bulk_storage" binding:"required"`
}

type TransactionByChainAndHash struct {
	ChainId int    `json:"chain_id" form:"chain_id" binding:"required"`
	Hash    string `json:"hash" form:"hash" binding:"required"`
}

type TransactionsByChainAndAddress struct {
	ChainId  int    `json:"chain_id" form:"chain_id" binding:"required"`
	Address  string `json:"address" form:"address" binding:"required"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"page_size"`
}

type SendMessageToTelegram struct {
	AuthKey string `json:"auth_key" form:"auth_key" binding:"required"`
	Message string `json:"message" form:"message" binding:"required"`
}

type RevokeTelegramKey struct {
	AuthKey string `json:"auth_key" form:"auth_key" binding:"required"`
}

type GetFreeCoin struct {
	Coin string `json:"coin" form:"coin" binding:"required"`
}
