package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	jellyfinAddress string
	jellyfinApiKey  string
	sessionsMetric  = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jellyfin_session_info",
			Help: "Information about Jellyfin sessions",
		},
		[]string{"UserName", "State", "Name", "Bitrate", "PlayMethod", "Substream", "DeviceName"},
	)
)

func init() {
	prometheus.MustRegister(sessionsMetric)
}

func main() {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	jellyfinApiKey = os.Getenv("JELLYFIN_APIKEY")

	if jellyfinAddress == "" || jellyfinApiKey == "" {
		log.Fatal("Please provide Jellyfin address and API key")
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	interval := 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	GetSessions()
	for {
		select {
		case <-ticker.C:
			GetSessions()
		}
	}
}
