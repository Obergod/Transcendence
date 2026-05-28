FROM node:alpine

WORKDIR /app

COPY frontend/package*.json ./

RUN npm install

COPY frontend .

# On récupère le WASM de l'image backend
COPY --from=backend-builder:latest /app/frontend/public/main.wasm ./public/main.wasm
COPY --from=backend-builder:latest /app/frontend/public/wasm_exec.js ./public/wasm_exec.js

RUN npm run build

COPY --from=backend-builder:latest /app/frontend/public/main.wasm /app/public/main.wasm
COPY --from=backend-builder:latest /app/frontend/public/wasm_exec.js /app/public/wasm_exec.js

CMD ["npm", "run", "dev", "--", "--host"]