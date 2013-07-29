package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// Handler for the fake MUC
func mucLogUIHandler(w http.ResponseWriter, r *http.Request) {
	// Check whether it's a request to add a key
	if r.FormValue("ID") != "" {
		ID, _ := strconv.Atoi(r.FormValue("ID"))
		addKeyToID(ID, r.FormValue("key"))
	}

	// Get all frames in the MUC table
	frames := getMucFrames()

	var rootTemplate = template.Must(template.ParseFiles("templates/muc.html"))
	rootTemplate.Execute(w, frames)
}
