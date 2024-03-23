#!/bin/bash

# Edit these to configure deployment
ONOMY_APP_IMAGE=onomy/onomyhub-hub-testnet:latest
ONOMY_APP_DOMAIN=app.testnet.onomy.io
ONOMY_APP_BCO_GQL_URL=https://api.thegraph.com/subgraphs/name/onomyprotocol/ograph-testnet
CMC_API_KEY=

sudo apt update && sudo apt upgrade -y
sudo apt install docker.io docker-compose tmux vim
sudo adduser ubuntu docker
sudo docker network create caddy

# Oracle Cloud Ubuntu Firewall Config
IPTABLES_CONFIG=/etc/iptables/rules.v4
if test -f "$IPTABLES_CONFIG"; then
  # Oracle Cloud Ubuntu Firewall Config
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 80 -j ACCEPT/' $IPTABLES_CONFIG
  sudo sed -i 's/22 -j ACCEPT/&\n-A INPUT -p tcp -m state --state NEW -m tcp --dport 443 -j ACCEPT/' $IPTABLES_CONFIG
  sudo iptables-restore < $IPTABLES_CONFIG
fi

cd ~/
echo "
version: '3'

services:
  caddy:
    image: lucaslorentz/caddy-docker-proxy:ci-alpine
    container_name: caddy
    ports:
      - 80:80
      - 443:443
    environment:
      - CADDY_INGRESS_NETWORKS=caddy
    networks:
      - caddy
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - caddy_data:/data
    restart: unless-stopped

  onomy-http-app:
    image: $ONOMY_APP_IMAGE
    restart: always
    container_name: onomy-http-app
    environment:
      - BCO_BACKEND_PORT=4444
      - BCO_GQL_BACKEND_URI=$ONOMY_APP_BCO_GQL_URL
      - CMC_API_KEY=$CMC_API_KEY
    networks:
      - caddy
    labels:
      caddy: $ONOMY_APP_DOMAIN
      caddy.reverse_proxy: "{{upstreams http 4444}}"

  watchtower:
    image: containrrr/watchtower:1.5.1
    restart: always
    container_name: watchtower
    environment:
      - WATCHTOWER_CLEANUP=true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --interval 30 --debug

  nodeexporter:
    image: prom/node-exporter:v1.3.1
    container_name: nodeexporter
    restart: always
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - '--path.procfs=/host/proc'
      - '--path.rootfs=/rootfs'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    ports:
      - '9100:9100'

networks:
  caddy:
    external: true

volumes:
  caddy_data: {}

" > ~/docker-compose.yml

sudo docker-compose up -d
