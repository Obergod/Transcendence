#!/usr/bin/env bash
set -e

if [ -z "$POSTGRES_PASSWORD" ]; then
  echo "POSTGRES_PASSWORD not set, generating one..."

  POSTGRES_PASSWORD=$(openssl rand -base64 24)

  echo "Generated password:"
  echo "$POSTGRES_PASSWORD"

  export POSTGRES_PASSWORD
fi

if [ ${#POSTGRES_PASSWORD} -lt 16 ]; then
  echo "POSTGRES_PASSWORD must be at least 16 characters."
  exit 1
fi

echo "Starting PostgreSQL..."