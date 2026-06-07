#!/bin/bash

mkdir -p grafana/ssl
chmod 755 grafana/ssl

if [ ! -f grafana/ssl/grafana.key ] || [ ! -f grafana/ssl/grafana.crt ]; then
    echo "Generating self-signed SSL certificate..."

    openssl req -new -newkey rsa:2048 -nodes -x509 -days 365 \
        -keyout grafana/ssl/grafana.key \
        -out grafana/ssl/grafana.crt \
        -subj "/CN=localhost" >/dev/null 2>&1

    chmod 644 grafana/ssl/grafana.crt
    chmod 600 grafana/ssl/grafana.key

    echo "Certificate generation complete."
fi