#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "Generating .env . . ."

ELASTIC_PASSWORD=$(openssl rand -base64 32)
KIBANA_PASSWORD=$(openssl rand -base64 32)

cat > .env <<EOF
ELASTIC_PASSWORD=$ELASTIC_PASSWORD
KIBANA_SYSTEM_PASSWORD=$KIBANA_PASSWORD
EOF

set -a
source .env
set +a

echo "Starting Elasticsearch . . ."

docker compose up -d elasticsearch

echo "Waiting for ElasticSearch . . ."
until curl -s -k -u elastic:$ELASTIC_PASSWORD https://localhost:9200 >/dev/null; do
	echo "Please wait . . ."
	sleep 2
done
echo "Done waiting!"

# ELASTIC-INIT

echo "Registering Snapshot repo . . .";
curl -s -k -u elastic:${ELASTIC_PASSWORD} -X PUT https://localhost:9200/_snapshot/s3-repo -H "Content-Type: application/json" -d @./src/elasticsearch/snapshot_repo.json;
echo "Repository registered.";


# ILM-INIT

echo "Registering Index Lifecycle Management policy . . .";
curl -s -k -u elastic:${ELASTIC_PASSWORD} -X PUT https://localhost:9200/_ilm/policy/logs-policy -H "Content-Type: application/json" -d @./src/elasticsearch/ilm_policy.json;
echo "ILM policy registered.";


# TEMPLATE-INIT

echo "Registering Template index. . .";
curl -s -k -u elastic:${ELASTIC_PASSWORD} -X PUT https://localhost:9200/_index_template/logs-template -H "Content-Type: application/json" -d @./src/elasticsearch/logs_template.json;
echo "Template registered.";

echo "Creating initial rollover index . . .";

curl -s -k -u elastic:${ELASTIC_PASSWORD} \
-X PUT https://localhost:9200/app-logs-000001 \
-H "Content-Type: application/json" \
-d '{
  "aliases": {
    "app-logs": {
      "is_write_index": true
    }
  }
}';

echo "Bootstrap index created.";


# SLM-INIT

echo "Registering Snapshot Lifecycle Management policy . . .";

curl -s -k -u elastic:${ELASTIC_PASSWORD} \
-X PUT https://localhost:9200/_slm/policy/daily-backups \
-H "Content-Type: application/json" \
-d '{
  "schedule": "0 0 2 * * ?",
  "name": "backup-{now/d}",
  "repository": "s3-repo",
  "config": {
    "indices": ["app-logs-*"],
    "include_global_state": false
  },
  "retention": {
    "expire_after": "14d",
    "min_count": 5,
    "max_count": 30
  }
}';

echo "SLM policy registered."

# REST



echo "Setting Kibana password . . ."

curl -s -k -X POST \
	-u elastic:$ELASTIC_PASSWORD \
	https://localhost:9200/_security/user/kibana_system/_password \
	-H "Content-Type: application/json" \
	-d "{\"password\":\"$KIBANA_PASSWORD\"}" >/dev/null

echo "Starting Kibana . . ."
docker compose up -d kibana

echo "Done."
exit