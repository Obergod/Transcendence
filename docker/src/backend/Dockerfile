FROM node

WORKDIR /app

COPY backend /app

RUN golang install

# COPY . .

RUN golang run build

RUN golang

EXPOSE 8081

CMD ["node", "dist/index.js"]