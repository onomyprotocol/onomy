* Create service file

```
vi /etc/systemd/system/your-service.service
```

Paste the content

```
[Unit]
Description=your-service

[Service]
User=root
WorkingDirectory=/root

ExecStart=/root/script.sh
Restart=always

[Install]
WantedBy=multi-user.target
```

* Reload the service files to include the new service.

```
systemctl daemon-reload
```

* Start your service

```
systemctl start your-service.service
```

* To check the status of your service

```
systemctl status your-service.service
```

* To enable your service on every reboot

```
systemctl enable your-service.service
```

* To disable your service on every reboot

```
systemctl disable your-service.service
```

* To get logs

```
journalctl -u your-service.service --no-pager
```
