##!/bin/bash

# Navigate to the directory containing the docker-compose file
cd "$(dirname "$0")"

# Stop any running containers
docker-compose down

# Build the services and force recreate the containers
docker-compose up --build --force-recreate
