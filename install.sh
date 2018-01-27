#!/usr/bin/env bash

# This script requires sshpass and parallel
#
passwd="bwcon"
export SSHPASS="${passwd}"


if [ -e "./workers.txt" ] ; then
	WORKERS=`cat ./workers.txt`
fi

if [ -z "${WORKERS}" ] ; then
	echo "Need some workers to install to!"
	exit 1
fi

dowork() {
	ipaddr=$1
	echo "----------- ${ipaddr} start"
	sshpass -e ssh -o StrictHostKeychecking=no root@$ipaddr "echo '
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
' > /etc/systemd/system/bw_exporter.service; sync; sync;"
	sshpass -e scp -o StrictHostKeychecking=no bw_exporter root@$ipaddr:/usr/bin
	sshpass -e ssh -o StrictHostKeychecking=no root@$ipaddr 'chmod ugo+x /usr/bin/bw_exporter; systemctl enable bw_exporter; systemctl start bw_exporter;'
	echo "----------- ${ipaddr} finish"
}

export -f dowork

parallel --no-notice dowork ::: ${WORKERS}
