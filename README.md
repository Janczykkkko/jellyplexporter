# jellyplexporter

Prometheus exporter for Plex and Jellyfin stream data.  

Main use case is visualising streams on a graph in grafana with the stream info as labels:  

>In Grafana (graph):
>
>select plex/jelly_session_info (both work the same)
>
>in Legend field add {{__name__}} / {{UserName}}({{DeviceName}}): {{Name}} ({{PlayMethod}}) subs: {{Substream}} bitrate: {{Bitrate}} Mbps
>
>Result label of each individual stream on a graph:
>
>jelly_session_info / USER(DEVICE): The Aviator (DirectPlay) subs: English - SUBRIP - External bitrate: 2.90243 Mbps

Made to be deployed as a container, requires at least one of two sets of variables:  
For Plex: PLEX_ADDRESS and PLEX_APIKEY  
For Jellyfin: JELLYFIN_ADDRESS and JELLYFIN_APIKEY  
*address in form >> http(s)://ip:port <<  

Images available here: https://hub.docker.com/repository/docker/januszadlo/jellyplexporter/general  
