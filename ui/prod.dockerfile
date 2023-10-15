FROM node:alpine AS builder

WORKDIR /app

COPY package.json ./
COPY package-lock.json ./
RUN npm ci

COPY . ./

RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.template /etc/nginx/templates/default.conf.template
COPY setup-env-config.sh /docker-entrypoint.d
RUN chmod +x /docker-entrypoint.d/setup-env-config.sh

ENV PORT 3000
ENV APP_DIR /usr/share/nginx/html