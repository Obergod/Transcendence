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

CMD ["./transcendance"]


# # =========================
# # 1. BACKEND + WASM BUILD
# # =========================
# FROM golang:1.25 AS backend-builder

# WORKDIR /app

# # Copy backend + shared code first
# COPY backend ./backend
# COPY internal ./internal
# COPY game ./game

# # If you REALLY want auto go.mod creation
# RUN if [ ! -f go.mod ]; then go mod init transcendance; fi

# ENV GOPROXY=https://proxy.golang.org,direct
# ENV GOSUMDB=off

# # Dependency download (cached)
# RUN --mount=type=cache,target=/go/pkg/mod \
#     go mod tidy

# # Build backend binary
# RUN --mount=type=cache,target=/root/.cache/go-build \
#     CGO_ENABLED=0 go build -o /out/transcendance ./backend

# # Build WASM output
# RUN mkdir -p /out/public && \
#     cd game && \
#     GOOS=js GOARCH=wasm go build -o /out/public/main.wasm main_wasm.go

# # Copy wasm runtime helper
# RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js /out/public/


# # =========================
# # 2. FRONTEND BUILD
# # =========================
# FROM node:24-bookworm AS frontend-builder

# WORKDIR /app

# COPY frontend/package*.json ./
# RUN npm install

# COPY frontend .

# # Inject WASM built by backend stage
# COPY --from=backend-builder /out/public/main.wasm ./public/main.wasm
# COPY --from=backend-builder /out/public/wasm_exec.js ./public/wasm_exec.js

# # Build frontend
# RUN npm run build


# # =========================
# # 3. RUNTIME IMAGE (SMALL)
# # =========================
# FROM node:24-bookworm

# WORKDIR /app

# # Only production artifacts
# COPY --from=frontend-builder /app ./

# # optional: backend binary if you want to run it in same container
# COPY --from=backend-builder /out/transcendance /usr/local/bin/transcendance

# EXPOSE 5173 8081

# CMD ["npm", "run", "dev", "--", "--host"]