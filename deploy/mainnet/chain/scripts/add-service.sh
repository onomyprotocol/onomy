#!/bin/bash
set -e

NAME=$1
START_SCRIPT=$2

echo "Adding ${NAME} service"

if [[ -z "${NAME}" ]]; then
  echo "NAME param is required"
  exit 1
fi

if test ! -f "$START_SCRIPT"; then
  echo "START_SCRIPT param is required, and the script must exist"
  exit 1
fi

echo "
[Unit]
Description=$NAME
After=network.target

[Service]
User=root
Environment=\"PATH=$PATH\"
Environment=\"DAEMON_HOME=$DAEMON_HOME\"
Environment=\"DAEMON_NAME=$DAEMON_NAME\"
Environment=\"DAEMON_RESTART_AFTER_UPGRADE=$DAEMON_RESTART_AFTER_UPGRADE\"
Environment=\"CHAIN_SCRIPTS=$PWD/\"
Environment=\"ETH_RPC_ADDRESS=$ETH_RPC_ADDRESS\"
Environment=\"VALIDATOR_AWS_KEYS_NAME=$VALIDATOR_AWS_KEYS_NAME\"
ExecStart=/bin/bash $START_SCRIPT
Restart=always
RestartSec=3
LimitAS=infinity
LimitRSS=infinity
LimitCORE=infinity
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
" > "/etc/systemd/system/$NAME.service"

systemctl daemon-reload
systemctl start "$NAME.service"
systemctl enable "$NAME.service"
systemctl status "$NAME.service" --no-pager
journalctl -u "$NAME.service" --no-pager