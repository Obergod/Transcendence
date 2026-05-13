FROM node

WORKDIR /app

COPY package*.json ./

RUN golang install

COPY . .

RUN golang run build

