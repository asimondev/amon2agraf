package amon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Account struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Server   string `json:"host"`
	Port     int    `json:"port"`
}

type Connect struct {
	Account
	config string
}

type JSONConfig struct {
	Account
	FileName   string `json:"file"`
	TableName  string `json:"table"`
	AppendRows bool   `json:"append"`
}

func parseJSONConfig(r io.Reader) (*JSONConfig, error) {
	data, _ := ioutil.ReadAll(r)
	var cfg JSONConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ReadJSONConfig(fn string) (*JSONConfig, error) {
	jsonFile, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	return parseJSONConfig(jsonFile)
}

func (a *Account) checkEnvVars() {
	if a.User == "" {
		if user := os.Getenv("AGRAF_MYSQL_USER"); user != "" {
			a.User = user
		}
	}

	if a.Password == "" {
		if pwd := os.Getenv("AGRAF_MYSQL_PASSWORD"); pwd != "" {
			a.Password = pwd
		}
	}

	if a.Database == "" {
		if db := os.Getenv("AGRAF_MYSQL_DATABASE"); db != "" {
			a.Database = db
		}
	}

	if a.Server == "" {
		if srv := os.Getenv("AGRAF_MYSQL_HOST"); srv != "" {
			a.Server = srv
		}
	}

	if a.Port == 0 {
		if port := os.Getenv("AGRAF_MYSQL_PORT"); port != "" {
			n, err := strconv.Atoi(port)
			if err != nil {
				log.Fatalf("AGRAF_MYSQL_PORT is not a number (%v)", err)
			}
			a.Port = n
		}
	}
}

func (acc *Account) String() string {
	return fmt.Sprintf("Account => User: '%s' Password: '%s' "+
		"Database: '%s' Server: '%s' Port: %d.\n",
		acc.User, acc.Password, acc.Server, acc.Database, acc.Port)
}

func NewDBAccount(u, p, d, s *string, n *int, cfg *JSONConfig) *Account {
	var acc Account = Account{User: *u, Password: *p, Server: *s, Database: *d, Port: *n}
	if cfg != nil {
		if acc.User == "" && len(cfg.User) > 0 {
			acc.User = cfg.User
		}

		if acc.Password == "" && len(cfg.Password) > 0 {
			acc.Password = cfg.Password
		}

		if acc.Database == "" && len(cfg.Database) > 0 {
			acc.Database = cfg.Database
		}

		if acc.Server == "" && len(cfg.Server) > 0 {
			acc.Server = cfg.Server
		}

		if acc.Port == 0 && cfg.Port > 0 {
			acc.Port = cfg.Port
		}
	}

	acc.checkEnvVars()

	return &acc
}

func CheckArgs(acc *Account) error {
	if acc.User == "" || acc.Password == "" || acc.Database == "" {
		return errors.New("user, password and database are mandatory")
	}

	if acc.Server == "" {
		acc.Server = "localhost"
	}

	if acc.Port == 0 {
		acc.Port = 3306
	}

	return nil
}
