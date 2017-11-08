package conf

import (
	"os"

	"github.com/BurntSushi/toml"
)

type conf struct {
	HTTPAddr    string `toml:"http_addr"`
	DbURI       string `toml:"db_uri"`
	SerialsURI  string `toml:"serials_uri"`
	SlackToken  string `toml:"slack_token"`
	KeyPemPath  string `toml:"key_pem_path"`
	CertPemPath string `toml:"cert_pem_path"`
}

var gc conf

func Load(path string) error {
	_, err := toml.DecodeFile(path, &gc)
	return err
}

func GetHTTPListenAddr() string {
	if gc.HTTPAddr == "" {
		return ":443"
	}
	return gc.HTTPAddr
}

func GetDatabaseURI() string {
	if gc.DbURI == "" {
		env := os.Getenv("DB_URI")
		if env != "" {
			return env
		}
	}
	return gc.DbURI
}

func GetSlackToken() string {
	if gc.SlackToken == "" {
		return os.Getenv("BOTTO_API_TOKEN")
	}
	return gc.SlackToken
}

func GetSerialsURI() string {
	return gc.SerialsURI
}

func GetKeyfilePath() string {
	return gc.KeyPemPath
}

func GetCertFilePath() string {
	return gc.CertPemPath
}
