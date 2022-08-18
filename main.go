package main

import (
	"fmt"
	"github.com/NurettinSelim/gdg-turkey-events-backend/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func getRoot(c *gin.Context) {
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
	c.JSON(http.StatusOK, data)
}

func getEvents(c *gin.Context) {
	queryType := database.QueryType(c.Query("queryType"))

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error occured while parsing %v", page)})
		return
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error occured while parsing %v", pageSize)})
		return
	}

	if !database.ValidQueries[queryType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong query parameter"})
		return
	}

	fsDatabase := database.FsDatabase{}

	if err = fsDatabase.Init(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fsDatabase.Close()

	events, err := fsDatabase.GetEvents(queryType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
func main() {
	r := gin.Default()
	r.GET("/", getRoot)

	api := r.Group("/api")
	{
		api.GET("/events", getEvents)
	}

	err := r.Run(":4000")
	if err != nil {
		log.Fatal(err)
		return
	}

}
