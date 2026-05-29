FROM golang:1.25 AS wasm-builder

WORKDIR /app

COPY backend ./backend
COPY internal ./internal
COPY game ./game

RUN [ -f go.mod ] || go mod init transcendance

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

RUN --mount=type=cache,target=/go/pkg/mod \
	GOPROXY=https://goproxy.cn,direct go mod tidy -v

RUN --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p frontend/public && \
    cd game && \
    GOARCH=wasm GOOS=js \
    go build -o ../frontend/public/main.wasm main_wasm.go

RUN cp $(go env GOROOT)/lib/wasm/wasm_exec.js frontend/public/

FROM node:24-bookworm

WORKDIR /app

COPY frontend/package*.json ./

RUN npm install

COPY frontend .

# On récupère le WASM de l'image builder
COPY --from=wasm-builder /app/frontend/public/main.wasm ./public/main.wasm
COPY --from=wasm-builder /app/frontend/public/wasm_exec.js ./public/wasm_exec.js

RUN npm run build

COPY --from=wasm-builder /app/frontend/public/main.wasm /app/public/main.wasm
COPY --from=wasm-builder /app/frontend/public/wasm_exec.js /app/public/wasm_exec.js

CMD ["npm", "run", "dev", "--", "--host"]
