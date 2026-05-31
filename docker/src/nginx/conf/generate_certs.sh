#!/bin/bash
SSL_DIR="/etc/nginx/ssl"

if [ ! -f "$SSL_DIR/cert.pem" ] || [ ! -f "$SSL_DIR/key.pem" ]; then
    echo "Generating self-signed certificates in $SSL_DIR..."

    # Generate private key
    openssl genrsa -out "$SSL_DIR/key.pem" 2048

    openssl req -new -key "$SSL_DIR/key.pem" -subj "/CN=localhost" -out "$SSL_DIR/cert.csr"
    openssl x509 -req -days 365 -in "$SSL_DIR/cert.csr" -signkey "$SSL_DIR/key.pem" -out "$SSL_DIR/cert.pem"

    rm "$SSL_DIR/cert.csr"

    chmod 600 "$SSL_DIR/key.pem"
    chmod 644 "$SSL_DIR/cert.pem"
fi

if [ -f "$SSL_DIR/cert.pem" ] && [ -f "$SSL_DIR/key.pem" ]; then
    echo "Certificates generated successfully!"
else
    echo "Failed to generate certificates!"
    exit 1
fi