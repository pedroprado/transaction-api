#!/bin/bash
docker-compose down

docker build -t apis:local .

docker-compose up -d --force-recreate apis postgres_db