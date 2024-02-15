package response

type AllInOneVerifyKeyResponse struct {
	Success bool   `json:"success"`
	Data    bool   `json:"data"`
	Msg     string `json:"msg"`
}
