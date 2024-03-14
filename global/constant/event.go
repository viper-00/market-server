package constant

type NotificationType string

var (
	INCOMING NotificationType = "incoming"
	OUTGOING NotificationType = "outgoing"
	SYSTEM   NotificationType = "system"
	OTHER    NotificationType = "other"
)

type ReadStatus int

var (
	DELETE ReadStatus = 0
	UNREAD ReadStatus = 1
	READ   ReadStatus = 2
)

type CoinIds string

var (
	IDS_Ethereum CoinIds = "ethereum"
	IDS_USDT     CoinIds = "tether"
	IDS_USDC     CoinIds = "usd-coin"
)

type EventType string

var (
	EVENT_CRYPTO   EventType = "Crypto"
	EVENT_BUSINESS EventType = "Business"
	EVENT_SCIENCE  EventType = "Science"
	EVENT_GAME     EventType = "Game"
)

type EventPlayType string

var (
	EVENT_PLAY_TWENTYSIX EventPlayType = "TwentySixLetters"
	EVENT_PLAY_THREE     EventPlayType = "ThreeLetters"
)

var EventPlayTypeValue []string

var (
	EVENT_PLAY_TWENTYSIX_VALUE = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "S", "Y", "Z"}
	EVENT_PLAY_THREE_VALUE     = []string{"A", "B", "C"}
)

var AllPlays = map[string][]string{
	string(EVENT_PLAY_TWENTYSIX): EVENT_PLAY_TWENTYSIX_VALUE,
	string(EVENT_PLAY_THREE):     EVENT_PLAY_THREE_VALUE,
}

var AllOrderTypes = map[uint]string{
	1: "buy",
	2: "sell",
}
