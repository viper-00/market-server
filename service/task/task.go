package task

import (
	MARKET_Client "market/utils/http"
)

var (
	client MARKET_Client.Client
)

func RunTask() {
	// go func() {
	//	RunApiKeyTestTask()
	// }()

	// go func() {
	// 	RunDailyReportTask()
	// }()

	// go func() {
	// 	RunPendingTxTask()
	// }()

	go func() {
		RunCoingeckoTask()
	}()
}
