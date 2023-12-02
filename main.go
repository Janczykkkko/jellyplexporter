package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/gin-gonic/gin"
)

var (
	pollingInterval = 30 * time.Second // Default polling interval
)

var (
	jellyfinAddress string
	apiKey          string
)

func getSessions() (int, error) {
	url := jellyfinAddress + "/Sessions?api_key=" + apiKey
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	log.Printf("API request to %s completed with status code: %d", jellyfinAddress, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return 0, err
	}
	result, err := jsonpath.Get("$[*].NowPlayingItem", jsonData)
	if err != nil {
		return 0, err
	}
	count := len(result.([]interface{}))
	return count, nil
}

func main() {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	apiKey = os.Getenv("API_KEY")
	pollingIntervalStr := os.Getenv("POLLING_INTERVAL")
	if pollingIntervalStr != "" {
		interval, err := strconv.Atoi(pollingIntervalStr)
		if err != nil {
			log.Fatalf("Invalid value for POLLING_INTERVAL: %s", err)
		}
		pollingInterval = time.Duration(interval) * time.Second
	}

	if jellyfinAddress == "" || apiKey == "" {
		log.Fatal("Please provide Jellyfin address and API key")
	}

	r := gin.Default()

	r.GET("/metrics", func(c *gin.Context) {
		count, err := getSessions()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error getting sessions count")
			return
		}
		c.String(http.StatusOK, "jellyfin_active_sessions_count %d", count)
	})
	log.Fatal(r.Run(":8080"))
}
