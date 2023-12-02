# jellyfin-exporter
Shitty jellyfin exporter written in go. 
Displays playback sessions count (couldn't find it anywhere on github). 
Two env vars needed: 
JELLYFIN_ADDRESS as in http://ip-or-hostname:port 
API_KEY api key lol :) 
Serves metrics at :8080/metrics
Confirmed working with jellyfin:10.8.10 inside a linuxserver image.
