FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY . .
RUN yarn build

FROM caddy:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY Caddyfile /etc/caddy/Caddyfile
EXPOSE 80
