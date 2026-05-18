FROM node

WORKDIR /app

COPY frontend/package*.json ./

RUN npm install

COPY frontend .

RUN npm run build

COPY --from=backend /app/frontend/public/main.wasm /app/public/main.wasm
COPY --from=backend /app/frontend/public/wasm_exec.js /app/public/wasm_exec.js


CMD ["npm", "run", "dev", "--", "--host"]