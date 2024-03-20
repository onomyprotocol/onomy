#!/bin/bash

## ulimits
echo 'fs.file-max = 65536' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
echo '* hard nofile 94000' | sudo tee -a /etc/security/limits.conf
echo '* soft nofile 94000' | sudo tee -a /etc/security/limits.conf
# The following line is necessary on RHEL/Oracle but not Ubuntu
# echo 'session required /lib/security/pam_limits.so' | sudo tee -a /etc/pam.d/login

