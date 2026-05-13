FROM node

WORKDIR /app

COPY frontend/package*.json ./

RUN npm install

COPY frontend .

RUN npm run build

# EXPOSE 8081

CMD ["npm", "run", "dev", "--", "--host"]

# FROM nginx

# # WORKDIR /app

# # RUN npm install -g serve

# COPY --from=builder /app/dist /usr/share/nginx/html

# EXPOSE 80
# CMD ["nginx", "-g", "daemon off;"]