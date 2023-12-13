package gatherers

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type PlexSessionMetric struct {
	UserName   string
	Name       string
	Bitrate    string
	PlayMethod string
	SubStream  string
	DeviceName string
	Value      int
}

// Get Plex data and parse it into a struct
func GetPlexSessions(plexAddress, plexApiKey string) (sessions PlexSessions, err error) {
	url := fmt.Sprintf(plexAddress + "/status/sessions?X-Plex-Token=" + plexApiKey)
	resp, err := http.Get(url)
	if err != nil {
		return PlexSessions{}, err
	}
	defer resp.Body.Close()
	log.Printf("API request to Plex at %s completed with status code: %d", plexAddress, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PlexSessions{}, err
	}
	err = xml.Unmarshal(body, &sessions)
	if err != nil {
		return PlexSessions{}, err
	}
	log.Println("Plex sessions scraped succesfully")
	return sessions, nil
}

// Ingest Plex data and assign metric per stream
func GetPlexMetrics(plexAddress, plexApiKey string) (metrics []PlexSessionMetric, err error) {
	sessions, err := GetPlexSessions(plexAddress, plexApiKey)
	if err != nil {
		return nil, err
	}
	for _, session := range sessions.Video {
		metric := PlexSessionMetric{
			UserName:   session.User.Title,
			Name:       session.Title,
			Bitrate:    GetPlexStreamBitrate(session),
			PlayMethod: session.Media.Part.Decision,
			SubStream:  GetPlexSubStream(session),
			DeviceName: session.Player.Device,
			Value:      1,
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

// convert bitrate
func GetPlexStreamBitrate(session PlexVideoSession) string {
	bitrateInt, err := strconv.Atoi(session.Media.Bitrate)
	if err != nil {
		log.Printf("Error processing plex stream bitrate: %s", err)
		return "Error"
	}
	return strconv.FormatFloat(float64(bitrateInt)/1000.0, 'f', -1, 64)
}

// stream type 3 is always the substream, need to find it
func GetPlexSubStream(session PlexVideoSession) (substream string) {
	substream = "None"
	for _, stream := range session.Media.Part.Stream {
		if stream.StreamType == "3" {
			substream = stream.ExtendedDisplayTitle
		}
	}
	return substream
}
