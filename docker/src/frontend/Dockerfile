FROM node

WORKDIR /app

COPY frontend /app

RUN npm install

RUN npm run build

# FROM nginx

# # WORKDIR /app

# # RUN npm install -g serve

# COPY --from=builder /app/dist /usr/share/nginx/html

# EXPOSE 80
# CMD ["nginx", "-g", "daemon off;"]