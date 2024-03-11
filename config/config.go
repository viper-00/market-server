package config

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type Server struct {
	System           System           `mapstructure:"system" json:"system" yaml:"system"`
	Redis            Redis            `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mysql            Mysql            `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Zap              Zap              `mapstructure:"zap" json:"zap" yaml:"zap"`
	Blockchain       Blockchain       `mapstructure:"blockchain" json:"blockchain" yaml:"blockchain"`
	BlockchainPlugin BlockchainPlugin `mapstructure:"blockchain-plugin" json:"blockchain-plugin" yaml:"blockchain-plugin"`
	Telegram         Telegram         `mapstructure:"telegram" json:"telegram" yaml:"telegram"`
	Wss              Wss              `mapstructure:"wss" json:"wss" yaml:"wss"`
	Smtp             Smtp             `mapstructure:"smtp" json:"smtp" yaml:"smtp"`
	Client           Client           `mapstructure:"client" json:"client" yaml:"client"`
	GeneralAccount   GeneralAccount   `mapstructure:"general-account" json:"general-account" yaml:"general-account"`
	Coingecko        Coingecko        `mapstructure:"coingecko" json:"coingecko" yaml:"coingecko"`
	File             File             `mapstructure:"file" json:"file" yaml:"file"`
}

type Coingecko struct {
	ApiKey string `mapstructure:"api-key" json:"api-key" yaml:"api-key"`
}

type GeneralAccount struct {
	Op Op `mapstructure:"op" json:"op" yaml:"op"`
}

type Op struct {
	PrivateKey     string `mapstructure:"private-key" json:"private-key" yaml:"private-key"`
	PublicKey      string `mapstructure:"public-key" json:"public-key" yaml:"public-key"`
	ReceiveAccount string `mapstructure:"receive-account" json:"receive-account" yaml:"receive-account"`
}

type Client struct {
	Url string `mapstructure:"url" json:"url" yaml:"url"`
}

type Smtp struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type Telegram struct {
	AllInOneInformBotLink        string `mapstructure:"allinone-inform-bot-link" json:"allinone-inform-bot-link" yaml:"allinone-inform-bot-link"`
	AllInOneInformBotToken       string `mapstructure:"allinone-inform-bot-token" json:"allinone-inform-bot-token" yaml:"allinone-inform-bot-token"`
	AllInOneInformChannelId      int64  `mapstructure:"allinone-inform-channel-id" json:"allinone-inform-channel-id" yaml:"allinone-inform-channel-id"`
	AllInOneNotificationBotLink  string `mapstructure:"allinone-notification-bot-link" json:"allinone-notification-bot-link" yaml:"allinone-notification-bot-link"`
	AllInOneNotificationBotToken string `mapstructure:"allinone-notification-bot-token" json:"allinone-notification-bot-token" yaml:"allinone-notification-bot-token"`
	AllInOneNotificationBotId    int64  `mapstructure:"allinone-notification-bot-id" json:"allinone-notification-bot-id" yaml:"allinone-notification-bot-id"`
}

type Wss struct {
	SecWssToken string `mapstructure:"sec-wss-token" json:"sec-wss-token" yaml:"sec-wss-token"`
}

type Blockchain struct {
	OpenSweepBlock bool `mapstructure:"open-sweep-block" json:"open-sweep-block" yaml:"open-sweep-block"`
	SweepMainnet   bool `mapstructure:"sweep-mainnet" json:"sweep-mainnet" yaml:"sweep-mainnet"`
	Ethereum       bool `mapstructure:"ethereum" json:"ethereum" yaml:"ethereum"`
	Bitcoin        bool `mapstructure:"bitcoin" json:"bitcoin" yaml:"bitcoin"`
	Tron           bool `mapstructure:"tron" json:"tron" yaml:"tron"`
	Bsc            bool `mapstructure:"bsc" json:"bsc" yaml:"bsc"`
	Litecoin       bool `mapstructure:"litecoin" json:"litecoin" yaml:"litecoin"`
	Op             bool `mapstructure:"op" json:"op" yaml:"op"`
	ArbitrumOne    bool `mapstructure:"arbitrum-one" json:"arbitrum-one" yaml:"arbitrum-one"`
	ArbitrumNova   bool `mapstructure:"arbitrum-nova" json:"arbitrum-nova" yaml:"arbitrum-nova"`
}

type BlockchainPlugin struct {
	Bitcoin  string `mapstructure:"bitcoin" json:"bitcoin" yaml:"bitcoin"`
	Litecoin string `mapstructure:"litecoin" json:"litecoin" yaml:"litecoin"`
}

type System struct {
	Env          string `mapstructure:"env" json:"env" yaml:"env"`
	DbType       string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr         int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	LimitCountIP int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP  int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	UseRedis     bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`
	UseInit      bool   `mapstructure:"use-init" json:"use-init" yaml:"use-init"`
	UseTask      bool   `mapstructure:"use-task" json:"use-task" yaml:"use-task"`
}

type File struct {
	ImageUrl        string `mapstructure:"image-url" json:"image-url" yaml:"image-url"`
	ImageClientPath string `mapstructure:"image-client-path" json:"image-client-path" yaml:"image-client-path"`
	ImageServerPath string `mapstructure:"image-server-path" json:"image-server-path" yaml:"image-server-path"`
	FileUrl         string `mapstructure:"file-url" json:"file-url" yaml:"file-url"`
	FileClientPath  string `mapstructure:"file-client-path" json:"file-client-path" yaml:"file-client-path"`
	FileServerPath  string `mapstructure:"file-server-path" json:"file-server-path" yaml:"file-server-path"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

type GeneralDB struct {
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Port         string `mapstructure:"port" json:"port" yaml:"port"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"`
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"`
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`
}

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func (m *Mysql) GetLogMode() string {
	return m.LogMode
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	MaxAge        int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`
}

func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder":
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder":
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder":
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder":
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
