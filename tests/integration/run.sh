#!/bin/bash

docker network create transaction_api_network
sh compose.sh

export PATH=$PATH:$PWD
echo "Waiting for all systems go..."
sleep 10

docker run --network transaction_api_network --rm --env API_URL=http://apis:8098/snapfi karate -t ~@ignore -T 1 /karate/cases/.
if [ $? -ne 0 ]; then
   echo "Error running integration tests"
   exit 1
fi

docker-compose stop
docker-compose down -v

rm -f docker-compose

exit 0
