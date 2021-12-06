#!/bin/bash
set -eu

sudo su -c "echo 'fs.file-max = 65536' >> /etc/sysctl.conf"
sysctl -p

sudo su -c "echo '* hard nofile 94000' >> /etc/security/limits.conf"
sudo su -c "echo '* soft nofile 94000' >> /etc/security/limits.conf"

sudo su -c "echo 'session session required /lib/security/pam_limits.so' >> /etc/pam.d/login"

echo "ulimit -n unlimited" >> ~/.bashrc

