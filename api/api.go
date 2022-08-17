package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type GDGApi struct{}

const GdgApiUrl = "https://gdg.community.dev/api/event/?fields=id,chapter,title,status,start_date,end_date,url&status=Published&start_date=%v&end_date=%v"

var countryMap = map[string]bool{
	"AF": true,
	"AM": true,
	"AZ": true,
	"KG": true,
	"KZ": true,
	"TM": true,
	"UZ": true,
	"MN": true,
	"TR": true,
}

type data struct {
	Links struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"links"`
	Count   int     `json:"count"`
	Results []Event `json:"results"`
}

type Event struct {
	Id        int       `json:"id" firestore:"id"`
	Title     string    `json:"title" firestore:"title"`
	Chapter   Chapter   `json:"chapter" firestore:"chapter"`
	StartDate time.Time `json:"start_date" firestore:"start_date"`
	EndDate   time.Time `json:"end_date" firestore:"end_date"`
	Url       string    `json:"url" firestore:"url"`
}

type Chapter struct {
	Country string `json:"country" firestore:"country"`
	Title   string `json:"title" firestore:"title"`
}

func (e *Event) String() string {
	return fmt.Sprintf("%d %s %s", e.Id, e.Title, e.Chapter)
}

func (e *Event) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":              e.Id,
		"title":           e.Title,
		"chapter_country": e.Chapter.Country,
		"chapter_title":   e.Chapter.Title,
		"start_date":      e.StartDate,
		"end_date":        e.EndDate,
		"url":             e.Url,
	}
}

func (c *Chapter) String() string {
	return fmt.Sprintf("Country: %s Title: %s", c.Country, c.Title)
}
func (g *GDGApi) GetEvents() []Event {
	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Now().AddDate(1, 0, 0)

	// Fetch Start
	fetchStart := time.Now()
	dataObj, err := getEventsFromApi(startTime, endTime, GdgApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	var events []Event
	events = append(events, dataObj.Results...)
	for dataObj.Links.Next != "" {
		dataObj, err = getEventsFromApi(time.Time{}, time.Time{}, dataObj.Links.Next)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, dataObj.Results...)
	}
	fetchEnd := time.Now()
	fmt.Printf("Fetch took %v\n", fetchEnd.Sub(fetchStart))

	i := 0
	for len(events) > i {
		if !countryMap[events[i].Chapter.Country] {
			events[i] = events[len(events)-1]
			events = events[:len(events)-1]
		} else {
			i++
		}
	}

	return events
}

func getEventsFromApi(startTime time.Time, endTime time.Time, url string) (data, error) {
	var resp *http.Response
	var err error

	if url == GdgApiUrl {
		resp, err = http.Get(fmt.Sprintf(GdgApiUrl, startTime.Format("2006-01-02"), endTime.Format("2006-01-02")))
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		return data{}, err
	}
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return data{}, readErr
	}

	var dataObj data

	jsonErr := json.Unmarshal(body, &dataObj)
	if jsonErr != nil {
		return data{}, jsonErr
	}

	return dataObj, nil
}
