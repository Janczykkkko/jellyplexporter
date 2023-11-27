# jellyfin-exporter
Shitty jellyfin exporter written in go. 
Displays sessions count. 
Two env vars needed: 
JELLYFIN_ADDRESS as in http://ip-or-hostname:port 
API_KEY api key lol :) 
Serves metrics at :8080/metrics
Confirmed working with jellyfin:10.8.10 inside a linuxserver image.
