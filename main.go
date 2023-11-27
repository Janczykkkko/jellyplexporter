package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	activeSessions = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jellyfin_active_sessions_count",
		Help: "Number of active sessions in Jellyfin",
	})
)

func init() {
	// Register the metric with Prometheus
	prometheus.MustRegister(activeSessions)
}

var (
	jellyfinAddress string
	apiKey          string
)

func getSessions() (int, error) {
	// Construct the URL with provided address and API key
	url := jellyfinAddress + "/Sessions?api_key=" + apiKey

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return 0, err
	}

	// Use jsonpath to extract the NowPlayingItem count
	result, err := jsonpath.Get("$[*].NowPlayingItem", jsonData)
	if err != nil {
		return 0, err
	}

	// Type assertion to extract the integer value
	count := len(result.([]interface{}))
	return count, nil
}

func main() {
	flag.StringVar(&jellyfinAddress, "address", "", "Jellyfin instance address")
	flag.StringVar(&apiKey, "apikey", "", "API key for Jellyfin")
	flag.Parse()

	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	apiKey = os.Getenv("API_KEY")

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
		activeSessions.Set(float64(count))
		c.String(http.StatusOK, "jellyfin_active_sessions_count %d", count)
	})

	go func() {
		log.Fatal(r.Run(":8080"))
	}()

	for {
		// Fetch session count
		count, err := getSessions()
		if err != nil {
			log.Printf("Error getting sessions count: %s", err)
		} else {
			// Update Prometheus metric with the new session count
			activeSessions.Set(float64(count))
		}

		// Sleep for 30 seconds before the next iteration
		time.Sleep(30 * time.Second)
	}
}
