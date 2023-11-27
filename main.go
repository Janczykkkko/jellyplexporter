package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		// Register the /metrics endpoint with Prometheus handler
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go func() {
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
	}()

	wg.Wait()
}
