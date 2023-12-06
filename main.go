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
		[]string{"UserName", "State", "Name", "PlayMethod", "Substream", "DeviceName"},
	)
	playingSessions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "jellyfin_playing_sessions_count",
			Help: "Number of sessions currently playing in Jellyfin",
		},
	)
)

func init() {
	prometheus.MustRegister(sessionsMetric)
	prometheus.MustRegister(playingSessions)
}

func main() {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	jellyfinApiKey = os.Getenv("API_KEY")

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

	for {
		select {
		case <-ticker.C:
			GetSessions()
		}
	}
}
