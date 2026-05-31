#!/bin/sh

mkdir -p /etc/nginx/ssl

chmod 700 /etc/nginx/ssl

/usr/local/bin/generate_certs.sh

exec "$@"