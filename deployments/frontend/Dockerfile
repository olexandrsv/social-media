ARG NODE_VERSION=25-alpine
ARG NGINX_VERSION=alpine3.22

FROM node:${NODE_VERSION} AS builder
WORKDIR /app/client
COPY client/package.json client/package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --omit=optional
COPY ./client/angular.json client/tsconfig*.json ./
COPY ./client/src ./src

RUN --mount=type=cache,target=/root/.npm \
    --mount=type=cache,target=/app/client/.angular/cache \
    npm run build -- --configuration=production

FROM nginxinc/nginx-unprivileged:${NGINX_VERSION} AS runner
USER nginx
COPY ./social-media/deployments/frontend/nginx.conf /etc/nginx/nginx.conf
COPY --chown=nginx:nginx --from=builder /app/client/dist/client/* /usr/share/nginx/html
EXPOSE 8080
ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
CMD ["-g", "daemon off;"]