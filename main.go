package main

import (
	"encoding/json"
	"fmt"
	"github.com/NurettinSelim/gdg-turkey-events-backend/database"
	"io"
	"log"
	"net/http"
	"strconv"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	type Author struct {
		Github string `json:"github"`
		Email  string `json:"email"`
	}
	data := struct {
		Author  Author `json:"author"`
		Version string `json:"version"`
	}{
		Author: Author{
			Github: "https://github.com/NurettinSelim",
			Email:  "nurettinselim03@gmail.com",
		},
		Version: "2.0.0",
	}
	jsonData, _ := json.Marshal(data)

	io.WriteString(w, string(jsonData))
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	queryType := database.QueryType(r.URL.Query().Get("queryType"))

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{"error":"Error occured while parsing %v"}`, page))
		return
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{"error":"Error occured while parsing %v"}`, pageSize))
		return

	}

	if !database.ValidQueries[queryType] {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Wrong query parameter"}`)
		return

	}

	fsDatabase := database.FsDatabase{}
	err = fsDatabase.Init()
	if err != nil {
		io.WriteString(w, "Error")
		return
	}
	defer fsDatabase.Close()

	events, err := fsDatabase.GetEvents(queryType, page, pageSize)

	if err != nil {
		io.WriteString(w, "Error")
		return
	}

	jsonData, _ := json.Marshal(events)

	io.WriteString(w, string(jsonData))
}
func main() {
	//gdgApi := api.GDGApi{}
	//events := gdgApi.GetEvents()

	//fsDatabase := database.FsDatabase{}
	//err := fsDatabase.Init()
	//defer fsDatabase.Close()
	//fmt.Println(fsDatabase.GetEventIds())
	//err = fsDatabase.SaveEvents(events)

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/api/events", getEvents)

	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
		return
	}

}
