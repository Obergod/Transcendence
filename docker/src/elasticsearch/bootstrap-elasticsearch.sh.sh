#!/bin/bash
set -e

ES_HOME=/usr/share/elasticsearch
ES_URL="http://localhost:9200"

echo "[bootstrap] Starting custom Elasticsearch bootstrap script" > /proc/1/fd/1

$ES_HOME/bin/elasticsearch &
ES_PID=$!

echo "[bootstrap] Waiting for HTTP endpoint..." >> /proc/1/fd/1
until curl -s "$ES_URL" >/dev/null; do
  sleep 2
done

echo "[bootstrap] waiting for cluster health..."
until curl -s "$ES_URL/_cluster/health" | grep -qE '"status":"(yellow|green)"'; do
  sleep 2
done

echo "[bootstrap] cluster is up"

if [ ! -f "/usr/share/elasticsearch/config/certs/elasticsearch.keystore" ]; then
  echo "[bootstrap] setting up security credentials..."

  curl -s -X POST -u "elastic:changeme" "$ES_URL/_security/user/elastic/_password" \
    -H "Content-Type: application/json" \
    -d "{\"password\": \"${ELASTIC_PASSWORD:-elastic123}\"}" >/dev/null

  curl -s -X POST -u "elastic:${ELASTIC_PASSWORD:-elastic123}" "$ES_URL/_security/user/kibana_system/_password" \
    -H "Content-Type: application/json" \
    -d "{\"password\": \"${KIBANA_PASSWORD:-kibana123}\"}" >/dev/null

  echo "[bootstrap] security credentials configured"
fi

wait $ES_PID