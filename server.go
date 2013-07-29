package main

import (
	"./mbus"
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

var sniffingTTY = flag.String("snifferTTY", "/dev/ttyUSB0", "sniffing device")
var sendingTTY = flag.String("senderTTY", "/dev/ttyUSB1", "sender device")
var DBUser = flag.String("DBUser", "root", "DB username")
var DBPass = flag.String("DBPass", "root", "DB password")
var DBName = flag.String("DBName", "capturedFrames", "DB name")
var DemoMode = flag.Bool("DemoMode", false, "Insert the sent data directly into the DB in case the sender is not working properly")

// Initialize the application
func main() {
	runtime.GOMAXPROCS(4)

	// Read config flags
	flag.Parse()

	// Initialize the database
	setupDB()

	// Initialize the sending device
	initializeSender()

	// Start the data reading from the serial device
	go readData()

	// Start the MUC webservice
	go mucService()

	// Start the sniffer webservice
	snifferService()
}

// Starts the MUC service on port 8080
func mucService() {
	http.HandleFunc("/webui", mucLogUIHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// Starts the sniffer webservice on port 80
func snifferService() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/send", sendFrameHandler)
	http.HandleFunc("/export", exportDump)

	http.HandleFunc("/import", importDump)
	http.HandleFunc("/truncate", truncateSQL)
	http.HandleFunc("/sendPartiallyHandler", sendPartiallyHandler)
	http.HandleFunc("/statusSniffer", statusSnifferHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var rootTemplate = template.Must(template.ParseFiles("templates/layout.html"))
	rootTemplate.Execute(w, snifferActive)
}

// Our socket handler for new sniffed data
// TODO: add a last ID
func socketHandler(c *websocket.Conn) {
	lastID := 0
	frames, ID := getNewFrames(lastID)
	json.NewEncoder(c).Encode(frames)
	lastID = ID

	for {
		frames, ID := getNewFrames(lastID)
		if len(frames) > 0 {
			json.NewEncoder(c).Encode(frames)
			lastID = ID
		}
	}
}

// Export a dump as JSON
func exportDump(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=\"dump-"+time.Now().String()+".json\"")
	frames, _ := getNewFrames(0)
	json.NewEncoder(w).Encode(frames)
}

// Import a dump
func importDump(w http.ResponseWriter, r *http.Request) {
	importedData, _, _ := r.FormFile("import")
	defer importedData.Close()

	var buf bytes.Buffer
	io.Copy(&buf, importedData)
	var frame []mbus.Frame
	json.Unmarshal(buf.Bytes(), &frame)
	for i := range frame {
		addNewFrame(frame[i].Value, frame[i].Time)
	}

	w.Header().Set("Refresh", "0; url=http://"+r.Host)
}

// Truncate SQL
func truncateSQL(w http.ResponseWriter, r *http.Request) {
	truncateTable()
	w.Header().Set("Refresh", "0; url=http://"+r.Host)
}
