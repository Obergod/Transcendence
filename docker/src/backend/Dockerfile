# FROM docker.io/library/golang:1.25 AS backend-builder

# WORKDIR /app

# COPY . .

# RUN [ -f go.mod ] || go mod init transcendance

# # COPY go.mod go.sum ./
# RUN go mod tidy

# RUN go mod download

# RUN go build -o transcendance ./backend

# # RUN mkdir -p frontend/public && cd game && GOARCH=wasm GOOS=js go build -o ../frontend/public/main.wasm main_wasm.go

# RUN mkdir -p frontend/public && \
#     cd game && \
#     GOARCH=wasm GOOS=js go build -x -v \
#     -o ../frontend/public/main.wasm main_wasm.go

# RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/

# # RUN cp /app/frontend/public/main.wasm /shared/main.wasm
# # RUN cp /usr/local/go/lib/wasm/wasm_exec.js /shared/wasm_exec.js

# # EXPOSE 5173

# CMD ["./transcendance"]








FROM golang:1.25 AS backend-builder

WORKDIR /app

RUN apt-get update && apt-get install -y openssl curl ca-certificates && \
    curl -LO https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64 && \
    mv mkcert-v1.4.4-linux-amd64 mkcert && \
    chmod +x mkcert && \
    mv mkcert /usr/local/bin/mkcert && \
    mkcert -install

RUN mkcert localhost 127.0.0.1 ::1

RUN mkdir -p /app/ssl && \
    cp /app/localhost+2.pem /app/ssl/localhost+2.pem && \
    cp /app/localhost+2-key.pem /app/ssl/localhost+2-key.pem

COPY backend ./backend
COPY internal ./internal
COPY game ./game

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

RUN mkdir -p /var/log/containers && \
    mkdir -p /app/logs

CMD ["./transcendance"]
