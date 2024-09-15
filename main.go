package main

import (
	"log"
	"net/http"
	"os"
	"time"

	gatherers "github.com/Janczykkkko/jellyplexgatherer"
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
	JellyOnlineMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jelly_online_info",
			Help: "Information about online Jelly users",
		},
		[]string{"UserName", "Device"},
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
		prometheus.MustRegister(JellyOnlineMetric)
	}
	if enablePlex {
		prometheus.MustRegister(PlexSessionMetric)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updateMetrics(enableJellyfin, enablePlex)
	}
}

func updateMetrics(enableJellyfin, enablePlex bool) {
	if enableJellyfin {
		JellyfinSessionsMetric.Reset() // reset to remove any bugged metrics
		sessionMetrics, err := gatherers.GetJellySessions(jellyfinAddress, jellyfinApiKey)
		if err != nil {
			log.Printf("Error getting Jellyfin session info: %s", err)
		}
		for _, metric := range sessionMetrics {
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

		JellyOnlineMetric.Reset()
		onlineMetrics, err := gatherers.GetOnlineUsers(jellyfinAddress, jellyfinApiKey, 100, 120)
		if err != nil {
			log.Printf("Error getting Jellyfin activity info: %s", err)
		}
		for _, metric := range onlineMetrics {
			activityLabels := prometheus.Labels{
				"UserName": metric.UserName,
				"Device":   metric.Device,
			}
			JellyOnlineMetric.With(activityLabels).Set(float64(1))
		}
	}
	if enablePlex {
		PlexSessionMetric.Reset() // reset to remove any bugged metrics
		metrics, err := gatherers.GetPlexSessions(plexAddress, plexApiKey)
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
