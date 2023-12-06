# jellyfin-exporter
Shitty jellyfin exporter written in go. 
Displays sessions count. 
Two env vars needed: 
JELLYFIN_ADDRESS as in http://ip-or-hostname:port 
JELLYFIN_APIKEY api key lol :) 
Serves metrics at :8080/metrics
Confirmed working with jellyfin:10.8.10 inside a linuxserver image.
https://hub.docker.com/repository/docker/januszadlo/jellyfin-exporter/general
Returns metric jellyfin_jellyfin_session_info{Bitrate="", DeviceName="", Name="", PlayMethod="Transcode", State="", Substream="", UserName=""} for each active stream.
