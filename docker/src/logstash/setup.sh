#!/bin/sh
set -e

mkdir -p /usr/share/logstash/pipeline

cat <<EOF > /usr/share/logstash/pipeline/logstash.conf
input {
  beats {
    port => 5044
  }
}

filter {
  mutate {
    add_field => {
      "environment" => "podman"
    }
  }
}

output {
  elasticsearch {
    hosts => ["http://elasticsearch:9200"]
    user => "elastic"
    password => "${ELASTIC_PASSWORD}"
  }

  stdout {
    codec => rubydebug
  }
}
EOF

exec /usr/local/bin/docker-entrypoint