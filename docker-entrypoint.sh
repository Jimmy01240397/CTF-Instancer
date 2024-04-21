#!/bin/sh

pwddir=$(pwd)
rm -f /var/run/docker.pid
dockerd &
sleep 3s
cd $CHALDIR
docker compose build
cd $pwddir
./ctfinstancer
