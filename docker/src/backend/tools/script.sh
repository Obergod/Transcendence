#!/bin/sh

go run /app/backend/main.go &

sleep 2

cd /app/frontend && npm run dev -- --host