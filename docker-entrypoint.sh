#!/bin/sh

pwddir=$(pwd)
rm -f /var/run/docker.pid
dockerd &
sleep 3s

for a in $(ls images)
do
    docker image load -i images/$a
done

cd $CHALDIR

for a in $(seq 0 1 $(($(yq '.services | length' docker-compose.yml) - 1)))
do
    imagename="$(yq ".services[.services | keys[$a]].image" docker-compose.yml)"
    if [ "$imagename" == "null" ]
    then
        echo "Please setup image name." 1>&2
        exit 1
    fi
    if ! docker images --format "{{.Repository}}:{{.Tag}}" | grep "$(echo "$imagename" | sed 's/\./\\./g')"
    then
        docker compose build "$(yq ".services | keys[$a]" docker-compose.yml)"
        docker image save "$imagename" > $pwddir/images/$(echo "$imagename" | awk -F: '{print $1}' | sed 's/\//_/g').tar
    fi
done

cd $pwddir

./ctfinstancer
