package amon

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Header struct {
	Host        string
	Format      string
	RAC         string
	CDB         string
	PDB         string
	InstName    string
	InstID      int
	User        string
	Interval    int
	StartTime   string `json:"Start Time"`
	EndTime     string
	Version     int
	AmonVersion string `json:"AMON Version"`
}

type FileHeader struct {
	AmonHeader Header
}

type Tail struct {
	EndTime string `json:"End Time"`
}

type FileTail struct {
	AmonTail Tail
}

func (hd Header) String() string {
	return fmt.Sprintf("AMON file details:\n"+
		" - Format: %s  Version: %d\n"+
		" - AMON Interval: %d\n"+
		" - AMON Report Start Time: %s  End Time: %s\n"+
		" - AMON Version: %s\n"+
		" - AMON DB User: %s\n"+
		" - RAC: %s  CDB: %s  PDB: %s\n"+
		" - Hostname: %s\n"+
		" - Instance Name: %s  Instance ID: %d\n",
		hd.Format, hd.Version, hd.Interval, hd.StartTime, hd.EndTime,
		hd.AmonVersion, hd.User, hd.RAC, hd.CDB, hd.PDB, hd.Host,
		hd.InstName, hd.InstID)
}

func Import(ver bool, fin string, tab string, app bool, mydb *sql.DB) error {
	header, err := readHeader(ver, fin)
	if err != nil {
		return err
	}

	header.EndTime = readTail(ver, fin)

	if header.Format == "amon_baseline_json" {
		return readBaseline(ver, fin, header, tab, app, mydb)
	}

	return fmt.Errorf("unknown AMON format type \"%s\"", header.Format)
}

func readHeader(ver bool, fin string) (*Header, error) {
	f, err := os.Open(fin)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	var fh FileHeader
	if err = dec.Decode(&fh); err != nil {
		return nil, fmt.Errorf("file header not found (%v)", err)
	}

	if ver {
		json, err := json.MarshalIndent(fh, "", "    ")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(json))
	}

	return &fh.AmonHeader, nil
}

// It's possible, that the file is truncated. So no errors are returned.
func readTail(ver bool, fin string) string {
	f, err := os.Open(fin)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	var fh FileTail
	if err = dec.Decode(&fh); err != nil {
		return ""
	}

	if ver {
		json, err := json.MarshalIndent(fh, "", "    ")
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(json))
	}

	return fh.AmonTail.EndTime
}
