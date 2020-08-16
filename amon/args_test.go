package amon

import (
	"strings"
	"testing"
)

func TestJsonConfigFile(t *testing.T) {
	var cfg *JSONConfig
	var err error
	if cfg, err = ReadJSONConfig("../cfg.json"); err != nil {
		t.Error(err)
	}

	if err != nil && cfg.User != "grafana" {
		t.Error("Wrong account user.")
	}
}

func TestJsonConfigUser(t *testing.T) {
	var cfg *JSONConfig
	var err error
	str := `{"User": "andrej"}`

	if cfg, err = parseJSONConfig(strings.NewReader(str)); err != nil {
		t.Error(err)
	}

	if err != nil && cfg.User != "andrej" {
		t.Error("Wrong account user.")
	}
}

func TestJsonConfigPassword(t *testing.T) {
	var cfg *JSONConfig
	var err error
	str := `{"User": "andrej", "Password": "haha"}`

	if cfg, err = parseJSONConfig(strings.NewReader(str)); err != nil {
		t.Error(err)
	}

	if err != nil && cfg.Password != "haha" {
		t.Error("Wrong account pasword.")
	}
}

func TestJsonConfigDatabase(t *testing.T) {
	var cfg *JSONConfig
	var err error
	str := `{"User": "andrej", "Password": "haha", "Database": "graf"}`

	if cfg, err = parseJSONConfig(strings.NewReader(str)); err != nil {
		t.Error(err)
	}

	if err != nil && cfg.Database != "graf" {
		t.Error("Wrong database.")
	}
}

func TestJsonConfigAmonFile(t *testing.T) {
	var cfg *JSONConfig
	var err error
	str := `{"User": "andrej", "Password": "haha", "File": "x.x"}`

	if cfg, err = parseJSONConfig(strings.NewReader(str)); err != nil {
		t.Error(err)
	}

	if err != nil && cfg.FileName != "x.x" {
		t.Error("Wrong file name.")
	}
}

func TestJsonConfigTableName(t *testing.T) {
	var cfg *JSONConfig
	var err error
	str := `{"User": "andrej", "Password": "haha", "Table": "tab"}`

	if cfg, err = parseJSONConfig(strings.NewReader(str)); err != nil {
		t.Error(err)
	}

	if err != nil && cfg.TableName != "tab" {
		t.Error("Wrong table name.")
	}
}

func TestJsonConfigAppendRows(t *testing.T) {
	var cfg *JSONConfig
	var err error
	str := `{"User": "andrej", "Password": "haha", "Table": "tab", "append": true}`

	if cfg, err = parseJSONConfig(strings.NewReader(str)); err != nil {
		t.Error(err)
	}

	if err != nil && cfg.AppendRows != true {
		t.Error("Wrong append value.")
	}
}
