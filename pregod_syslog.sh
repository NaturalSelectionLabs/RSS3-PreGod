#!/bin/sh
 
parentDir="/etc/rsyslog.d"
fileName="pregod_syslog.conf"
dirAndName="$parentDir/$fileName"
if [ ! -d "$parentDir" ];then
mkdir $parentDir
fi

cd $parentDir
 
if [ ! -f "$dirAndName" ];then
cat>$dirAndName<<EOF
local0.* action(type="omfwd"
        protocol="tcp"
        target="$PROMTAIL_IP"
        port="$PROMTAIL_PORT"
        Template="RSYSLOG_SyslogProtocol23Format"
        TCP_Framing="octet-counted")
EOF
echo $fileName "Created"
else
echo $fileName "Already Exists"
fi
