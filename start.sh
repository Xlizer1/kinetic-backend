#!/bin/bash

# Get Docker container IPs
DB_IP=$(docker inspect kinetic-db --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' 2>/dev/null)
REDIS_IP=$(docker inspect kinetic-redis --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' 2>/dev/null)

# Fallback to defaults if containers not running
DB_IP=${DB_IP:-172.21.0.2}
REDIS_IP=${REDIS_IP:-172.21.0.3}

cd /opt/kite/kinetic-backend

export SERVER_PORT=8080
export DB_HOST=$DB_IP
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD='KineticSecureDB2024!'
export DB_NAME=kinetic
export JWT_SECRET=35097284c884ca13c365760c93d3dbf0079158d4b8beb02716e4db6d7cc58105940fbe0bfe26d4771b3d6b30e86317422a1d18c7d63ff396ac0664f92fb5c62c
export REDIS_HOST=$REDIS_IP
export REDIS_PORT=6379
export REDIS_PASSWORD=bf5f78b06d0bf5558bedad96655a2851724779b0d06f8e3d0d6beabec491e2e4
export LIVEKIT_API_KEY=APIRES8Xpqw437u
export LIVEKIT_API_SECRET=a7VN5rSfbuXHvPsSU8YZ0t67FyEE9kWBVcgOGska6hG
export LIVEKIT_SERVER_URL=wss://kinetic-66jvnjn7.livekit.cloud
export GIN_MODE=release

exec ./kinetic-server