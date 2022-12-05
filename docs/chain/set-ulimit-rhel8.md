# Set correct ulimits for rhel8

Run the following commands to set correct ulimits.

```
echo 'fs.file-max = 65536' >> /etc/sysctl.conf
sysctl -p

echo '* hard nofile 94000' >> /etc/security/limits.conf
echo '* soft nofile 94000' >> /etc/security/limits.conf

echo 'session session required /lib/security/pam_limits.so' >> /etc/pam.d/login
```
