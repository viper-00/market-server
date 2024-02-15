package tatum

type TatumGetLitecoinInfo struct {
	Chain         string  `json:"chain"`
	Blocks        int     `json:"blocks"`
	Headers       int     `json:"headers"`
	Bestblockhash string  `json:"bestblockhash"`
	Difficulty    float64 `json:"difficulty"`
}

type TatumGetLitecoinBlock struct {
	Hash       string            `json:"hash"`
	Height     int               `json:"height"`
	Depth      int               `json:"depth"`
	Version    int               `json:"version"`
	PrevBlock  string            `json:"prevBlock"`
	MerkleRoot string            `json:"merkleRoot"`
	Time       int               `json:"time"`
	Bits       int               `json:"bits"`
	Nonce      int               `json:"nonce"`
	Txs        []TatumLitecoinTx `json:"txs"`
}

type TatumLitecoinTx struct {
	Hash        string             `json:"hash"`
	Hex         string             `json:"hex"`
	WitnessHash string             `json:"witnessHash"`
	Fee         string             `json:"fee"`
	Rate        string             `json:"rate"`
	Mtime       int                `json:"mtime"`
	BlockNumber int                `json:"blockNumber"`
	Block       string             `json:"block"`
	Time        int                `json:"time"`
	Index       int                `json:"index"`
	Version     int                `json:"version"`
	Locktime    int                `json:"locktime"`
	Inputs      []LitecoinTxInput  `json:"inputs"`
	Outputs     []LitecoinTxOutput `json:"outputs"`
}

type LitecoinTxInput struct {
	Prevout  LitecoinTxInputPrevout `json:"prevout"`
	Script   string                 `json:"script"`
	Witness  string                 `json:"witness"`
	Sequence int                    `json:"sequence"`
	Coin     LitecoinTxInputCoin    `json:"coin"`
}

type LitecoinTxInputPrevout struct {
	Hash  string `json:"hash"`
	Index int    `json:"index"`
}

type LitecoinTxInputCoin struct {
	Version  int    `json:"version"`
	Height   int    `json:"height"`
	Value    string `json:"value"`
	Script   string `json:"script"`
	Address  string `json:"address"`
	Coinbase bool   `json:"coinbase"`
}

type LitecoinTxOutput struct {
	Value   string `json:"value"`
	Script  string `json:"script"`
	Address string `json:"address"`
}
