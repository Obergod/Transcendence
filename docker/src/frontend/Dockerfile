FROM node

WORKDIR /app

COPY frontend /app

RUN npm install

RUN npm run build

COPY --from=backend /app/frontend/public/main.wasm /app/public/main.wasm
COPY --from=backend /app/frontend/public/wasm_exec.js /app/public/wasm_exec.js


CMD ["npm", "run", "dev", "--", "--host"]
