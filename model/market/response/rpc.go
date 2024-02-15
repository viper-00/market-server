package response

type RPCBlockInfo struct {
	JsonRpc string             `json:"jsonrpc"`
	Id      int                `json:"id"`
	Result  RPCBlockInfoResult `json:"result"`
}

type RPCBlockInfoResult struct {
	BaseFeePerGas    string   `json:"baseFeePerGas"`
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Transactions     []string `json:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}

type RPCBlockDetail struct {
	JsonRpc string               `json:"jsonrpc"`
	Id      int                  `json:"id"`
	Result  RPCBlockDetailResult `json:"result"`
}

type RPCBlockDetailResult struct {
	BaseFeePerGas    string           `json:"baseFeePerGas"`
	Difficulty       string           `json:"difficulty"`
	ExtraData        string           `json:"extraData"`
	GasLimit         string           `json:"gasLimit"`
	GasUsed          string           `json:"gasUsed"`
	Hash             string           `json:"hash"`
	LogsBloom        string           `json:"logsBloom"`
	Miner            string           `json:"miner"`
	MixHash          string           `json:"mixHash"`
	Nonce            string           `json:"nonce"`
	Number           string           `json:"number"`
	ParentHash       string           `json:"parentHash"`
	ReceiptsRoot     string           `json:"receiptsRoot"`
	Sha3Uncles       string           `json:"sha3Uncles"`
	Size             string           `json:"size"`
	StateRoot        string           `json:"stateRoot"`
	Timestamp        string           `json:"timestamp"`
	TotalDifficulty  string           `json:"totalDifficulty"`
	Transactions     []RPCTransaction `json:"transactions"`
	TransactionsRoot string           `json:"transactionsRoot"`
	Uncles           []string         `json:"uncles"`
}

type RPCTransactionDetail struct {
	JsonRpc string         `json:"jsonrpc"`
	Id      int            `json:"id"`
	Result  RPCTransaction `json:"result"`
}

type RPCTransaction struct {
	BlockHash            string `json:"blockHash"`
	BlockNumber          string `json:"blockNumber"`
	Hash                 string `json:"hash"`
	From                 string `json:"from"`
	Gas                  string `json:"gas"`
	GasPrice             string `json:"gasPrice"`
	Input                string `json:"input"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	Nonce                string `json:"nonce"`
	R                    string `json:"r"`
	S                    string `json:"s"`
	SourceHash           string `json:"sourceHash"`
	To                   string `json:"to"`
	TransactionIndex     string `json:"transactionIndex"`
	Type                 string `json:"type"`
	V                    string `json:"v"`
	Value                string `json:"value"`
	ChainId              string `json:"chainId"`
}

type RPCGeneral struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}
