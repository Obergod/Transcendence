FROM golang:1.25 AS backend-builder

WORKDIR /app

COPY backend /app

RUN [ -f go.mod ] || go mod init transcendance

# COPY go.mod go.sum ./
RUN go mod tidy

RUN go mod download

RUN golang run build

# RUN mkdir -p frontend/public && cd game && GOARCH=wasm GOOS=js go build -o ../frontend/public/main.wasm main_wasm.go

RUN mkdir -p frontend/public && \
    cd game && \
    GOARCH=wasm GOOS=js go build -x -v \
    -o ../frontend/public/main.wasm main_wasm.go

RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/
# EXPOSE 5173

CMD ["./transcendance"]
