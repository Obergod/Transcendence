FROM node:24-alpine

WORKDIR /app

RUN apk update && apk add openssl curl ca-certificates &&  \
    wget https://github.com/FiloSottile/mkcert/releases/download/v1.4.4/mkcert-v1.4.4-linux-amd64 && \
    mv mkcert-v1.4.4-linux-amd64 mkcert && \
    chmod +x mkcert && \
    mv mkcert /usr/local/bin/mkcert && \
    mkcert -install

# RUN mkcert localhost 127.0.0.1 ::1

COPY frontend/package*.json ./

RUN npm install

COPY frontend .

COPY --from=backend-builder /app/ssl/localhost+2-key.pem /app/ssl/localhost+2-key.pem
COPY --from=backend-builder /app/ssl/localhost+2.pem /app/ssl/localhost+2.pem

RUN npm run build

COPY --from=backend-builder /app/frontend/public/main.wasm /app/public/main.wasm
COPY --from=backend-builder /app/frontend/public/wasm_exec.js /app/public/wasm_exec.js

CMD ["npm", "run", "dev", "--", "--host"]
