[![Build Status](https://travis-ci.org/lookfirst/bw_exporter.svg?branch=master)](https://travis-ci.org/lookfirst/bw_exporter)

# BW Exporter

[Prometheus.io](https://prometheus.io/) exporter for the cgminer binary that is shipped with the BW.com L21 miner. It currently exports a limited set of data. PR's welcome!

Thanks to [HyperBit.io](https://hyperbitshop.io) for sponsoring this project.

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

The [releases tab](https://github.com/lookfirst/bw_exporter/releases) has `master` binaries cross compiled for ARM suitable for running on the miner. These are built automatically on [Travis](https://travis-ci.org/lookfirst/bw_exporter).

Download the latest release and copy the `bw_exporter` binary to `/usr/bin`

```
chmod ugo+x bw_exporter
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

Open your browser to `http://MINER_IP:4030/metrics`

### Prometheus configuration

`prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'bw_exporter'
    file_sd_configs:
      - files:
        - 'bw_exporter.json'
```

`bw_exporter.json`:
```json
[{
	"targets": ["MINER_IP:4030"]
}]
```

The json configuration is reloaded every time it is modified.