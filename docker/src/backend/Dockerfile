FROM docker.io/library/golang:1.25 AS backend-builder

WORKDIR /app

COPY . .

RUN [ -f go.mod ] || go mod init transcendance

# COPY go.mod go.sum ./
RUN go mod tidy

RUN go mod download

RUN go build -o transcendance ./backend

# RUN mkdir -p frontend/public && cd game && GOARCH=wasm GOOS=js go build -o ../frontend/public/main.wasm main_wasm.go

RUN mkdir -p frontend/public && \
    cd game && \
    GOARCH=wasm GOOS=js go build -x -v \
    -o ../frontend/public/main.wasm main_wasm.go

RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/

# RUN cp /app/frontend/public/main.wasm /shared/main.wasm
# RUN cp /usr/local/go/lib/wasm/wasm_exec.js /shared/wasm_exec.js

# EXPOSE 5173

CMD ["./transcendance"]
