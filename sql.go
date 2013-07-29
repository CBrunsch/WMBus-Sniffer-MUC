package main

import (
	"./mbus"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
)

var db *sql.DB

// Initializes the DB
func setupDB() {
	db, _ = sql.Open("mysql", "root:root@/capturedFrames?parseTime=true&loc=Local")
}

// Indicates whether new data is available
var newData int

// Add a frame to the database
func addFrameToDB(frame *mbus.Frame) {
	stmt, err := db.Prepare("INSERT INTO sniffedFrames (value) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	newData++
	stmt.Exec(frame.Value)

	// Try to update the frame in the muc table
	stmt, err = db.Prepare("UPDATE muc SET `value`=?,`timestamp`=? WHERE address=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(frame.Value, time.Now(), frame.Identification())

	// Try to insert the frame in the muc table
	// this won't work if there is already one as the column is defined unique
	stmt, err = db.Prepare("INSERT INTO muc (`value`, `address`, `key`) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	var key string
	if frame.Address() == "440000570C37" {
		key = "CAFEBABE123456789ABCDEF0CAFEBABE"
	} else {
		key = ""
	}
	defer stmt.Close()
	stmt.Exec(frame.Value, frame.Identification(), key)
}

// Add key to MUC ID
func addKeyToID(ID int, key string) {
	stmt, err := db.Prepare("UPDATE muc SET `key`=? WHERE `ID`=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(key, ID)
}

// Get all muc frames
func getMucFrames() []mbus.Frame {
	rows, err := db.Query("SELECT `ID`, `value`, `timestamp`, `key` FROM muc ORDER BY timestamp DESC")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch rows
	var frames []mbus.Frame

	for rows.Next() {
		var ID int
		var value string
		var timestamp time.Time
		var key string
		rows.Scan(&ID, &value, &timestamp, &key)
		frame, _ := mbus.NewFrame(value)
		frame.Time = timestamp
		frame.ID = ID
		frame.Key = key
		frames = append(frames, *frame)
	}
	return frames
}

// Get new frames
func getNewFrames(lastFrame int) ([]mbus.Frame, int) {
	rows, err := db.Query("SELECT * FROM sniffedFrames WHERE id > " + strconv.Itoa(lastFrame))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch rows
	lastRow := lastFrame
	var frames []mbus.Frame

	for rows.Next() {
		var ID int
		var value string
		var timestamp time.Time
		rows.Scan(&ID, &value, &timestamp)
		if ID != 0 {
			frame, _ := mbus.NewFrame(value)
			hexbyte, _ :=  hex.DecodeString(fmt.Sprintf("%s",value))
			frame.Hexified = fmt.Sprintf("% X",fmt.Sprintf("%s", hexbyte))
			frame.Time = timestamp
			frame.ID = ID
			frames = append(frames, *frame)
			lastRow = ID
		}
	}
	return frames, lastRow
}

// Truncates the SQL table
func truncateTable() {
	db.Query("TRUNCATE TABLE sniffedFrames")
	db.Query("TRUNCATE TABLE muc")
}

// Add new frame to SQL 
func addNewFrame(value string, timestamp time.Time) {
	stmt, err := db.Prepare("INSERT INTO sniffedFrames (value, timestamp) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	stmt.Exec(value, timestamp)
}
