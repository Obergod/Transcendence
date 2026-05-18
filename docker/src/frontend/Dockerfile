# FROM node

# WORKDIR /app

# COPY frontend/package*.json ./

# RUN npm install

# COPY frontend .

# RUN npm run build

# # EXPOSE 8081

# CMD ["npm", "run", "dev", "--", "--host"]

# # FROM nginx

# # # WORKDIR /app

# # # RUN npm install -g serve

# # COPY --from=builder /app/dist /usr/share/nginx/html

# # EXPOSE 80
# # CMD ["nginx", "-g", "daemon off;"]


FROM node:20

WORKDIR /app

COPY . .

# install frontend deps
RUN cd frontend && npm install

# install Go temporarily
RUN apt-get update && apt-get install -y golang

# build wasm
RUN mkdir -p /app/frontend/public && \
    cd /app/game && \
    GOOS=js GOARCH=wasm go build \
    -o /app/frontend/public/main.wasm main_wasm.go

# copy wasm_exec.js
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" \
    /app/frontend/public/

WORKDIR /app/frontend

CMD ["npm", "run", "dev", "--", "--host"]