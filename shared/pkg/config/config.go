package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/util"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

// Golbal Config

type ProtocolStruct struct {
	Version string `koanf:"version"`
}

type HubServerStruct struct {
	RunMode      string        `koanf:"run_mode"`
	HttpPort     int           `koanf:"http_port"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
}

type RedisStruct struct {
	Addr     string `koanf:"addr"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`

	TLSEnabled bool `koanf:"tls_enabled"`
}

type PostgresStruct struct {
	DSN             string        `koanf:"dsn"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
}

type MongoStruct struct {
	URI         string `koanf:"uri"`
	DB          string `koanf:"db"`
	MaxPoolSize int    `koanf:"max_pool_size"`
	MinPoolSize int    `koanf:"min_pool_size"`
}

type BrokerStruct struct {
	Addr string `koanf:"addr"`
}

/*
 * File module：You can simply make the log print to a file
 * example:
 * Type = "file"
 * FilePath = "./log/app.log"
 *
 * Syslog module: You can make the log print to a syslog server
 * * You need to download rsyslog.
 * * you need to use syslog.sh under 'scripts' to generate the configuration file pregod_syslog.conf under /etc/rsyslog.d
 * * The default configuration of facility is 0, which is related to the generated configuration file, see local*
 */
type LoggerOutputConfig struct {
	Type     string `koanf:"type"`     // available values: `stdout`, `file`, `syslog`
	Filepath string `koanf:"filepath"` // only for file
	Facility int    `koanf:"facility"` // only for syslog, available values: `1 - 7` to `LOG_LOCAL0-LOG_LOCAL7`
}

type LoggerStruct struct {
	PrefixTag string `koanf:"prefix_tag"`
	Engine    string `koanf:"engine"`   // available values: `zap`
	Level     string `koanf:"level"`    // available values: `debug`, `info`, `warn`, `error`, `panic`, `fatal`
	Encoding  string `koanf:"encoding"` // available values: `json`, `console`

	Output []LoggerOutputConfig `koanf:"output"`
}

type NetWorkStruct struct {
	Proxy string `koanf:"proxy"`
}

// Indexer Struct Config

type MoralisStruct struct {
	ApiKey string `koanf:"api_key"`
}

type EtherScanStruct struct {
	ApiKey string `koanf:"api_key"`
}

type PolygonScanStruct struct {
	ApiKey string `koanf:"api_key"`
}

type ArbitrumStruct struct {
	ApiKey string `koanf:"api_key"`
}

type MiscStruct struct {
	UserAgent string `koanf:"user_agent"`
}

//nolint:tagliatelle // format is required by Jike API
type JikeStruct struct {
	AreaCode          string `koanf:"area_code" json:"areaCode"`
	MobilePhoneNumber string `koanf:"mobile_phone_number" json:"mobilePhoneNumber"`
	Password          string `koanf:"password" json:"password"`
	AppVersion        string `koanf:"app_version" json:"appVersion"`
}
type TwitterStruct struct {
	Tokens []string `koanf:"break_down_tokens"`
}
type IndexerServerStruct struct {
	RunMode      string        `koanf:"run_mode"`
	HttpPort     int           `koanf:"http_port"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
}

type IndexerStruct struct {
	Server IndexerServerStruct `koanf:"server"`

	Misc        MiscStruct        `koanf:"misc"`
	Jike        JikeStruct        `koanf:"jike"`
	Moralis     MoralisStruct     `koanf:"moralis"`
	EtherScan   EtherScanStruct   `koanf:"etherscan"`
	PolygonScan PolygonScanStruct `koanf:"polygonscan"`
	Twitter     TwitterStruct     `koanf:"twitter"`
	Aribtrum    ArbitrumStruct    `koanf:"arbitrum"`
}

type ConfigStruct struct {
	Protocol  ProtocolStruct  `koanf:"protocol"`
	HubServer HubServerStruct `koanf:"hub_server"`
	Redis     RedisStruct     `koanf:"redis"`
	Postgres  PostgresStruct  `koanf:"postgres"`
	Mongo     MongoStruct     `koanf:"mongo"`
	Broker    BrokerStruct    `koanf:"broker"`
	Logger    LoggerStruct    `koanf:"logger"`
	Network   NetWorkStruct   `koanf:"network"`
	Indexer   IndexerStruct   `koanf:"indexer"`
}

var (
	Config = &ConfigStruct{}

	k = koanf.New(".")
)

func Setup() error {
	// Read user config
	fp, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := k.Load(file.Provider(fp), json.Parser()); err != nil {
		return err
	}

	if err := k.Unmarshal("", Config); err != nil {
		return err
	}

	Config.HubServer.ReadTimeout = Config.HubServer.ReadTimeout * time.Second
	Config.HubServer.WriteTimeout = Config.HubServer.WriteTimeout * time.Second

	Config.Postgres.ConnMaxIdleTime = Config.Postgres.ConnMaxIdleTime * time.Second
	Config.Postgres.ConnMaxLifetime = Config.Postgres.ConnMaxLifetime * time.Second

	Config.Indexer.Server.ReadTimeout = Config.Indexer.Server.ReadTimeout * time.Second
	Config.Indexer.Server.WriteTimeout = Config.Indexer.Server.WriteTimeout * time.Second

	return nil
}

// Gets config file path.
// Config files are located at `config/config.*.json`.
// The wildcard part is specified with an env var `CONFIG_ENV`.
// The default `CONFIG_ENV` is `dev`; that is, the default config file is `config/config.dev.json`.
func getConfigFilePath() (string, error) {
	ce := os.Getenv("CONFIG_ENV")
	if ce == "" {
		os.Setenv("CONFIG_ENV", "dev")

		ce = "dev"
	}

	dirname, err := util.Dirname()
	if err != nil {
		return "", err
	}

	fp := filepath.Join(dirname, "..", "..", "..", "config", "config."+ce+".json")

	return fp, nil
}
