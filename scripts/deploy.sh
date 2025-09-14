#!/bin/sh
set -e

# Accept positional arguments
BT_MAILGUNPASSWORD="$1"
BT_ENV="$2"

# Optionally check for required args
if [ -z "$BT_MAILGUNPASSWORD" ] || [ -z "$BT_ENV" ]; then
  echo "Usage: $0 <BT_MAILGUNPASSWORD> <BT_ENV>"
  exit 1
fi

# Export for docker-compose
export BT_MAILGUNPASSWORD
export BT_ENV

cd /home/prinzpiuz/BloTils
git fetch origin
git checkout release
git pull origin release

rm -rf ../BloTils_Data/config.json
cp -r config.json ../BloTils_Data/.

docker compose -f docker/docker-compose.yml pull
docker compose -f docker/docker-compose.yml up -d --remove-orphans
docker image prune -f
docker compose -f docker/docker-compose.yml ps
