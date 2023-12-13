package main

import (
	"jellyplexporter/gatherers"
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
	plexAddress       string
	plexApiKey        string
	PlexSessionMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "plex_session_info",
			Help: "Information about Plex sessions",
		},
		[]string{"UserName", "Name", "Bitrate", "PlayMethod", "Substream", "DeviceName"},
	)
)

func main() {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	jellyfinApiKey = os.Getenv("JELLYFIN_APIKEY")
	plexAddress = os.Getenv("PLEX_ADDRESS")
	plexApiKey = os.Getenv("PLEX_APIKEY")
	enableJellyfin, enablePlex := checkEnvs()

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	if enableJellyfin {
		prometheus.MustRegister(JellyfinSessionsMetric)
	}
	if enablePlex {
		prometheus.MustRegister(PlexSessionMetric)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updateMetrics(enableJellyfin, enablePlex)
	}
}

func updateMetrics(enableJellyfin, enablePlex bool) {
	if enableJellyfin {
		JellyfinSessionsMetric.Reset() // reset to remove any bugged metrics
		metrics, err := gatherers.GetJellyMetrics(jellyfinAddress, jellyfinApiKey)
		if err != nil {
			log.Printf("Error getting Jellyfin session info: %s", err)
		}
		for _, metric := range metrics {
			sessionLabels := prometheus.Labels{
				"UserName":   metric.UserName,
				"Name":       metric.Name,
				"Bitrate":    metric.Bitrate,
				"PlayMethod": metric.PlayMethod,
				"Substream":  metric.SubStream,
				"DeviceName": metric.DeviceName,
			}
			JellyfinSessionsMetric.With(sessionLabels).Set(float64(1))
		}
	}
	if enablePlex {
		PlexSessionMetric.Reset() // reset to remove any bugged metrics
		metrics, err := gatherers.GetPlexMetrics(plexAddress, plexApiKey)
		if err != nil {
			log.Printf("Error getting Plex session info: %s", err)
		}
		for _, metric := range metrics {
			sessionLabels := prometheus.Labels{
				"UserName":   metric.UserName,
				"Name":       metric.Name,
				"Bitrate":    metric.Bitrate,
				"PlayMethod": metric.PlayMethod,
				"Substream":  metric.SubStream,
				"DeviceName": metric.DeviceName,
			}
			PlexSessionMetric.With(sessionLabels).Set(float64(1))
		}
	}
}

func checkEnvs() (enableJellyfin, enablePlex bool) {
	enableJellyfin = true
	enablePlex = true
	if jellyfinAddress == "" || jellyfinApiKey == "" {
		log.Println("Jellyfin address and API key, Jellyfin metrics will not be exported")
		enableJellyfin = false
	}
	if plexAddress == "" || plexApiKey == "" {
		log.Println("Plex address and API key, Plex metrics will not be exported")
		enablePlex = false
	}
	return enableJellyfin, enablePlex
}
