# jellyplexporter

Prometheus exporter for Plex and Jellyfin.  
Main use case is visualising streams on a graph in grafana with the stream info as labels:
>{USER(DEVICE): The Aviator (DirectPlay) subs: English - SUBRIP - External bitrate: 2.90243 Mbps} 1

Made to be deployed as a container, requires at least one of two sets of variables:  
For Plex: PLEX_ADDRESS and PLEX_APIKEY  
For Jellyfin: JELLYFIN_ADDRESS and JELLYFIN_APIKEY  
*address in form >> http(s)://ip:port <<  

Images available here: https://hub.docker.com/repository/docker/januszadlo/jellyplexporter/general
