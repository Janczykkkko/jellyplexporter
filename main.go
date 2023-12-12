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
	jellyfinAddress        string
	jellyfinApiKey         string
	JellyfinSessionsMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jellyfin_session_info",
			Help: "Information about Jellyfin sessions",
		},
		[]string{"UserName", "Name", "Bitrate", "PlayMethod", "Substream", "DeviceName"},
	)
)

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
	prometheus.MustRegister(JellyfinSessionsMetric)
	interval := 30 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	GetJellyfinSessions()
	for {
		select {
		case <-ticker.C:
			GetJellyfinSessions()
		}
	}
}
