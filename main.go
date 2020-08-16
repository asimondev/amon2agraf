package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"amon2agraf/amon"
)

func main() {
	var app = flag.Bool("append", false, "append rows to existing table")
	var f = flag.String("file", "", "AMON data file")
	var d = flag.String("database", "", "database name")
	var n = flag.Int("port", 0, "port number")
	var p = flag.String("password", "", "database user password")
	var s = flag.String("host", "", "host name")
	var tab = flag.String("table", "", "table name")
	var u = flag.String("user", "", "database user")
	var ver = flag.Bool("version", false, "version")
	var v = flag.Bool("verbose", false, "Verbose output")
	var cf = flag.String("config", "", "JSON config file")

	flag.Parse()

	if *ver {
		fmt.Println("amon2agraf version: 1.0.0")
		return
	}

	var err error
	var cfg *amon.JSONConfig
	if len(*cf) > 0 {
		cfg, err = amon.ReadJSONConfig(*cf)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	}

	acc := amon.NewDBAccount(u, p, d, s, n, cfg)

	fileName, err := setFileName(f, cfg)
	if err != nil {
		log.Fatal(err)
	}

	tableName := setTableName(tab, cfg)
	appendRows := setAppendRows(*app, cfg)

	if err = amon.CheckArgs(acc); err != nil {
		log.Fatal("Error: ", err)
	}

	if *v {
		fmt.Println(acc)
		fmt.Printf("File name: %s\nTable name: %s\nAppend rows: %t\n",
			fileName, tableName, appendRows)
	}

	var mydb *sql.DB

	mydb, err = amon.OpenMySQL(acc, *v)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer mydb.Close()

	err = amon.Import(*v, *f, tableName, appendRows, mydb)
	if err != nil {
		log.Fatal(err)
	}
}

func setFileName(fn *string, cfg *amon.JSONConfig) (string, error) {
	if *fn != "" {
		return *fn, nil
	}

	if cfg != nil && cfg.FileName != "" {
		return cfg.FileName, nil
	}

	return "", fmt.Errorf("Error: AMON file name is missing")
}

// Check for the table name. If not found, returns "amon_baseline".
func setTableName(tab *string, cfg *amon.JSONConfig) string {
	if *tab != "" {
		return *tab
	}

	if cfg != nil && cfg.TableName != "" {
		return cfg.TableName
	}

	return "amon_baseline_tab"
}

// Check for the append rows flag.
func setAppendRows(app bool, cfg *amon.JSONConfig) bool {
	if app {
		return true
	}

	if cfg != nil && cfg.AppendRows {
		return true
	}

	return false
}
