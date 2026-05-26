FROM golang:1.25-alpine AS backend-builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOTMPDIR=/app/build-tmp
RUN mkdir -p /app/build-tmp

COPY . .

RUN [ -f go.mod ] || go mod init transcendance

RUN go mod tidy
RUN go mod download

RUN go build -o transcendance ./backend

RUN mkdir -p frontend/public && \
    cd game && \
    GOARCH=wasm GOOS=js go build -o ../frontend/public/main.wasm main_wasm.go

RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/

CMD ["./transcendance"]