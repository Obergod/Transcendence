FROM node AS builder

WORKDIR /app

# 1. copy dependency files FIRST
COPY package*.json ./

# 2. install dependencies
RUN npm install

# 3. copy the rest of the app
COPY . .

# 4. build
RUN npm run build


FROM nginx

# WORKDIR /app

# RUN npm install -g serve

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]