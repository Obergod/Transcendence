#!/bin/bash

sudo touch /etc/containers/nodocker

SSL_DIR="./src/elasticsearch/ssl"
CONFIG_DIR=".."
PODMAN_DIR="/usr/share/elasticsearch/config/ssl"
ELASTIC_PASSWORD="your_secure_password_here"

mkdir -p "$SSL_DIR"
mkdir -p "$CONFIG_DIR"

openssl genrsa -out "$SSL_DIR/ca.key" 2048
openssl req -new -x509 -days 365 -key "$SSL_DIR/ca.key" \
    -subj "/CN=Elasticsearch Transport CA" \
    -out "$SSL_DIR/ca.crt"

openssl pkcs12 -export -noprompt \
    -in "$SSL_DIR/ca.crt" \
    -inkey "$SSL_DIR/ca.key" \
    -out "$SSL_DIR/ca.p12" \
    -password pass:changeme

keytool -genkey -noprompt -dname "CN=elasticsearch-http" \
    -alias elasticsearch-http \
    -keyalg RSA \
    -keysize 2048 \
    -keystore "$SSL_DIR/elasticsearch.p12" \
    -storepass changeme \
    -storetype PKCS12

keytool -genkey -noprompt -dname "CN=elasticsearch-transport" \
    -alias elasticsearch-transport \
    -keyalg RSA \
    -keysize 2048 \
    -keystore "$SSL_DIR/transport.p12" \
    -storepass changeme \
    -storetype PKCS12

# Note: The above keytool commands already create PKCS12 files
# But if you need to create a separate http keystore:
# keytool -genkey -noprompt -dname "CN=elasticsearch-http" \
#     -alias elasticsearch-http \
#     -keyalg RSA \
#     -keysize 2048 \
#     -keystore "$SSL_DIR/elasticsearch.p12" \
#     -storepass changeme \
#     -storetype PKCS12

podman run -d --name elasticsearch \
    --restart=on-failure \
    --network=appnet \
    -p 9200:9200 \
    -e cluster.name=docker-cluster \
    -e network.host=0.0.0.0 \
    -e discovery.type=single-node \
    -e bootstrap.memory_lock=true \
    -e xpack.security.enabled=true \
    -e ELASTIC_PASSWORD="$ELASTIC_PASSWORD" \
    -e xpack.security.http.ssl.enabled=true \
    -e xpack.security.http.ssl.keystore.path="$PODMAN_DIR/elasticsearch.p12" \
    -e xpack.security.http.ssl.truststore.path="$PODMAN_DIR/ca.p12" \
    -e xpack.security.http.ssl.keystore.password=changeme \
    -e xpack.security.http.ssl.truststore.password=changeme \
    -e xpack.security.transport.ssl.enabled=true \
    -e xpack.security.transport.ssl.keystore.path="$PODMAN_DIR/transport.p12" \
    -e xpack.security.transport.ssl.truststore.path="$PODMAN_DIR/ca.p12" \
    -e xpack.security.transport.ssl.keystore.password=changeme \
    -e xpack.security.transport.ssl.truststore.password=changeme \
    -e xpack.security.transport.ssl.verification_mode=certificate \
    -v esdata:/usr/share/elasticsearch/data \
    -v "$CONFIG_DIR/src/elasticsearch/ssl:$PODMAN_DIR" \
    -ulimit memlock=-1 \
    docker.elastic.co/elasticsearch/elasticsearch:9.4.0