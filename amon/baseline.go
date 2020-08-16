package amon

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Baseline struct {
	ID        int64
	Timestamp string
	InstID    int
	InstName  string
	ConID     int
	StatID    int
	StatName  string
	Value     float64
}

type BaselineData struct {
	Data []Baseline
}

func readBaseline(ver bool, fin string, hd *Header, tab string, app bool,
	mydb *sql.DB) error {

	f, err := os.Open(fin)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	var data BaselineData
	if err = dec.Decode(&data); err != nil {
		return fmt.Errorf("data not found (%v)", err)
	}

	if ver {
		json, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(json))
	}

	if len(data.Data) == 0 {
		return errors.New("no data found")
	}

	err = data.Insert(mydb, hd, tab, app)
	if err != nil {
		return err
	}

	err = data.CreateView(mydb, tab)
	if err != nil {
		return err
	}

	fmt.Println(hd)
	fmt.Printf("AMON file details inserted into the table %s.\n", AmonTable)
	fmt.Printf("%d AMON baseline rows inserted into the table %s.\n",
		len(data.Data), tab)
	fmt.Println("View amon_baseline was re-created.")

	return nil
}

func (bl *BaselineData) Insert(mydb *sql.DB, hd *Header, tab string, app bool) error {
	id, err := hd.Insert(mydb, tab)
	if err != nil {
		return err
	}

	if !app {
		truncateTable(mydb, tab)
	}

	if err := bl.InsertData(mydb, tab, id); err != nil {
		return bl.CheckInsert(mydb, tab, id, err)
	}

	return nil
}

func (bl *BaselineData) CreateInsertData(mydb *sql.DB, tab string, id int64) error {
	stmt := "create table " + tab +
		" (id mediumint not null, amon_id mediumint not null, " +
		"time_stamp datetime not null, " +
		"inst_id int not null, inst_name varchar(128) not null," +
		"con_id int not null, " +
		"stat_id int not null, " +
		"stat_name varchar(128) not null, " +
		"stat_value double not null)"

	if _, err := mydb.Exec(stmt); err != nil {
		return err
	}

	return bl.InsertData(mydb, tab, id)
}

func (bl *BaselineData) CheckInsert(mydb *sql.DB, tab string, id int64, err error) error {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number != 1146 {
			return err
		}

		return bl.CreateInsertData(mydb, tab, id)
	}

	return err
}

func (bl *BaselineData) InsertData(mydb *sql.DB, tab string, id int64) error {
	tx, err := mydb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into " + tab +
		" (id, amon_id, time_stamp, inst_id, inst_name, con_id, stat_id," +
		"stat_name, stat_value) " +
		"values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range bl.Data {
		res, err := stmt.Exec(row.ID, id, row.Timestamp, row.InstID, row.InstName,
			row.ConID, row.StatID, row.StatName, row.Value)
		if err != nil {
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows != 1 {
			return errors.New("no rows were inserted")
		}
	}

	return tx.Commit()
}

func (bl *BaselineData) CreateView(mydb *sql.DB, tab string) error {
	stmt := "create or replace view amon_baseline " +
		"as select * from " + tab

	if _, err := mydb.Exec(stmt); err != nil {
		return err
	}

	return nil
}
