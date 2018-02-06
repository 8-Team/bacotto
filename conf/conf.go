package conf

import (
	"os"

	"github.com/BurntSushi/toml"
)

type API struct {
	HTTPAddr    string `toml:"http_addr"`
	UseHTTPS    bool   `toml:"use_https"`
	VerifyMsg   bool   `toml:"verify_msg"`
	KeyPemPath  string `toml:"key_pem_path"`
	CertPemPath string `toml:"cert_pem_path"`
}

type Bot struct {
	SlackToken string `toml:"slack_token"`
}

type DB struct {
	DBURI      string `toml:"db_uri"`
	SerialsURI string `toml:"serials_uri"`
	FixtureURI string `toml:"fixture_uri"`
}

type ERP struct {
	UseMock      bool   `toml:"use_mock"`
	MockDataPath string `toml:"mock_data_path"`
}

type Log struct {
	Debug bool `toml:"debug"`
}

type conf struct {
	API API `toml:"api"`
	Bot Bot `toml:"bot"`
	DB  DB  `toml:"db"`
	ERP ERP `toml:"erp"`
	Log Log `toml:"log"`
}

var gc conf

func Load(path string) error {
	_, err := toml.DecodeFile(path, &gc)
	return err
}

func GetHTTPListenAddr() string {
	if gc.API.HTTPAddr == "" {
		return ":443"
	}
	return gc.API.HTTPAddr
}

func UseHTTPS() bool {
	return gc.API.UseHTTPS
}

func VerifyMsg() bool {
	return gc.API.VerifyMsg
}

func GetDatabaseURI() string {
	if gc.DB.DBURI == "" {
		env := os.Getenv("DB_URI")
		if env != "" {
			return env
		}
	}
	return gc.DB.DBURI
}

func GetSlackToken() string {
	if gc.Bot.SlackToken == "" {
		return os.Getenv("BOTTO_API_TOKEN")
	}
	return gc.Bot.SlackToken
}

func GetSerialsURI() string {
	return gc.DB.SerialsURI
}

func GetFixtureURI() string {
	return gc.DB.FixtureURI
}

func GetKeyfilePath() string {
	return gc.API.KeyPemPath
}

func GetCertFilePath() string {
	return gc.API.CertPemPath
}

func UseMockERP() bool {
	return gc.ERP.UseMock
}

func MockERPDataPath() string {
	return gc.ERP.MockDataPath
}

func DebugLogLevel() bool {
	return gc.Log.Debug
}
