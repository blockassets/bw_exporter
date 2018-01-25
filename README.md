# BW Exporter

This is a Prometheus.io exporter for the cgminer binary that is shipped with the BW.com L21 miner. It currently only exports data for the chipstat command, but will export more data in the near future.

### Usage (defaults):

``
./bw_exporter -port 4030 -cghost 127.0.0.1 -cgport 4028 -cgtimeout 5s 
``

### Setup

Install [dep](https://github.com/golang/dep) and the dependencies...

`make dep`

### Build binary for arm

`make arm`

### Install onto miner

Copy the `bw_exporter` binary to `/usr/bin`

```
scp bw_exporter root@MINER_IP:/usr/bin
```

Create `/etc/systemd/system/bw_exporter.service`

```
ssh root@MINER_IP "echo '
[Unit]
Description=bw_exporter
After=init.service

[Service]
Type=simple
ExecStart=/usr/bin/bw_exporter
Restart=always
RestartSec=4s
StandardOutput=journal+console

[Install]
WantedBy=multi-user.target
' > /etc/systemd/system/bw_exporter.service"
```

Enable the service:

```
ssh root@MINER_IP "systemctl enable bw_exporter; systemctl start bw_exporter"
```

### Test install on miner

Open your browser to http://MINER_IP:4030/metrics
