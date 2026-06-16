# FROM golang:1.25 AS builder

# WORKDIR /app

# RUN apt-get update && apt-get install -y openssl curl ca-certificates && \
#     curl -LO https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64 && \
#     mv mkcert-v1.4.4-linux-amd64 mkcert && \
#     chmod +x mkcert && \
#     mv mkcert /usr/local/bin/mkcert && \
#     mkcert -install

# RUN mkcert localhost 127.0.0.1 ::1

# RUN mkdir -p /app/ssl && \
#     cp /app/localhost+2.pem /app/ssl/localhost+2.pem && \
#     cp /app/localhost+2-key.pem /app/ssl/localhost+2-key.pem

# COPY backend ./backend
# COPY internal ./internal
# COPY game ./game

# RUN if [ ! -f go.mod ]; then \
#         go mod init transcendance; \
#     fi

# ENV GOPROXY=https://proxy.golang.org,direct
# ENV GOSUMDB=off

# RUN --mount=type=cache,target=/go/pkg/mod \
# 	GOPROXY=https://goproxy.cn,direct go mod tidy -v

# RUN --mount=type=cache,target=/root/.cache/go-build \
#     --mount=type=cache,target=/go/pkg/mod \
#     CGO_ENABLED=0 \
#     go build -o transcendance ./backend

# RUN mkdir -p frontend/public && \
#     cd game && \
#     GOARCH=wasm GOOS=js go build -o ../frontend/public/main.wasm main_wasm.go

# RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/wasm_exec.js
# # RUN --mount=type=cache,target=/root/.cache/go-build \
# #     mkdir -p frontend/public && \
# #     cd game && \
# #     GOARCH=wasm GOOS=js \
# #     go build -o ../frontend/public/main.wasm main_wasm.go

# # RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/

# FROM node:24-alpine AS frontend-builder

# WORKDIR /app

# COPY frontend/package*.json ./
# RUN npm install

# COPY frontend .
# RUN npm run build

# FROM alpine

# WORKDIR /app

# RUN apk add --no-cache ca-certificates

# COPY --from=builder /app/transcendance /app/
# COPY --from=builder /app/ssl /app/ssl
# COPY --from=builder /app/frontend/public/main.wasm /app/public/
# COPY --from=builder /app/frontend/public/wasm_exec.js /app/public/

# COPY --from=frontend-builder /app/dist /app/public

# RUN mkdir -p /app/logs

# # RUN mkdir -p /var/log/containers && \
# #     mkdir -p /app/logs

# CMD ["./transcendance"]


FROM golang:1.25 AS backend-builder

WORKDIR /app

RUN apt-get update && apt-get install -y openssl curl ca-certificates npm && \
    curl -LO https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64 && \
    mv mkcert-v1.4.4-linux-amd64 mkcert && \
    chmod +x mkcert && \
    mv mkcert /usr/local/bin/mkcert && \
    mkcert -install

RUN mkcert localhost 127.0.0.1 ::1

RUN mkdir -p /app/ssl && \
    cp /app/localhost+2.pem /app/ssl/localhost+2.pem && \
    cp /app/localhost+2-key.pem /app/ssl/localhost+2-key.pem

COPY . .

RUN if [ ! -f go.mod ]; then \
        go mod init transcendance; \
    fi

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod tidy

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 \
    go build -o transcendance ./backend

RUN --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p frontend/public && \
    cd game && \
    GOARCH=wasm GOOS=js \
    go build -o ../frontend/public/main.wasm main_wasm.go

RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/

RUN cd frontend && npm install && npm run build

RUN mkdir -p /var/log/containers && \
    mkdir -p /app/logs

COPY docker/src/backend/tools/script.sh /script.sh
RUN chmod +x /script.sh

ENTRYPOINT ["/script.sh"]
