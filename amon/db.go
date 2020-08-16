package amon

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	AmonTable = "agraf_amon"
)

func OpenMySQL(acc *Account, verbose bool) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", acc.User, acc.Password,
		acc.Server, acc.Port, acc.Database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func truncateTable(mydb *sql.DB, tab string) error {
	stmt := "truncate table " + tab

	_, err := mydb.Exec(stmt)

	return err
}

func (hd *Header) CreateInsertHeader(mydb *sql.DB, tab string) (int64, error) {
	stmt := "create table " + AmonTable +
		" (id mediumint not null auto_increment, " +
		"host varchar(128) not null, " + "rac varchar(3) not null, " +
		"inst_id int not null, inst_name varchar(128) not null," +
		"cdb varchar(3) not null, pdb varchar(128)," +
		"amon_user varchar(128) not null, " +
		"amon_interval int not null," +
		"amon_format varchar(128) not null, " +
		"amon_version varchar(32) not null," +
		"start_time datetime not null, end_time datetime, " +
		"table_name varchar(64) not null, " +
		"version int not null," +
		"PRIMARY KEY (id))"

	if _, err := mydb.Exec(stmt); err != nil {
		return 0, err
	}

	return insertHeader(mydb, hd, tab)
}

func insertHeader(mydb *sql.DB, hd *Header, tab string) (int64, error) {
	tx, err := mydb.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into " + AmonTable +
		" (host, rac, inst_id, inst_name, cdb, pdb," +
		"amon_user, amon_interval, amon_format, amon_version, " +
		"start_time, end_time, table_name, version) " +
		"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var endTime sql.NullString
	if hd.EndTime != "" {
		endTime.String = hd.EndTime
		endTime.Valid = true
	}

	res, err := stmt.Exec(hd.Host, hd.RAC, hd.InstID, hd.InstName,
		hd.CDB, hd.PDB, hd.User, hd.Interval, hd.Format, hd.AmonVersion,
		hd.StartTime, endTime, tab, hd.Version)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rows != 1 {
		return 0, errors.New("no rows were inserted")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can not get last ID (%v)", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (hd *Header) CheckInsert(mydb *sql.DB, tab string, err error) (int64, error) {
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number != 1146 {
			return 0, err
		}

		return hd.CreateInsertHeader(mydb, tab)
	}

	return 0, err
}

// Insert AMON header. Create a table AmonTable, if it does not exist.
// The first return value is the inserted ID.
func (hd *Header) Insert(mydb *sql.DB, tab string) (int64, error) {
	id, err := insertHeader(mydb, hd, tab)
	if err != nil {
		return hd.CheckInsert(mydb, tab, err)
	}

	return id, nil
}
