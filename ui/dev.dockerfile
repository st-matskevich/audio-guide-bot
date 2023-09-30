FROM node:alpine AS builder

WORKDIR /app

COPY package.json ./
COPY package-lock.json ./
RUN npm ci

COPY setup-env-config.js ./

ENTRYPOINT ["npm", "start"]