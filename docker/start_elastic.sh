#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "Generating env..."

ELASTIC_PASSWORD=$(openssl rand -base64 32)
KIBANA_PASSWORD=$(openssl rand -base64 32)

cat > .env <<EOF
ELASTIC_PASSWORD=$ELASTIC_PASSWORD
KIBANA_SYSTEM_PASSWORD=$KIBANA_PASSWORD
EOF

echo "Starting Elasticsearch..."

docker compose up -d elasticsearch

echo "Waiting for ES..."
until curl -s -u elastic:$ELASTIC_PASSWORD http://localhost:9200 >/dev/null; do
  sleep 2
done

echo "Setting kibana_system password..."

curl -s -X POST \
  -u elastic:$ELASTIC_PASSWORD \
  http://localhost:9200/_security/user/kibana_system/_password \
  -H "Content-Type: application/json" \
  -d "{\"password\":\"$KIBANA_PASSWORD\"}" >/dev/null

echo "Starting Kibana..."
docker compose up -d kibana

echo "Done."