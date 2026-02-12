FROM golang:1.25-alpine AS builder
ARG SERVICE
WORKDIR /project
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY ./go.mod ./go.sum ./
#RUN go mod download
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
WORKDIR /project/cmd/${SERVICE}
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o app .


FROM alpine:latest
ARG SERVICE PORT
ARG APP_PATH="internal/common/app"
WORKDIR /cmd
# RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S user \
 && adduser -S -G user -H -s /sbin/nologin user
COPY --from=builder --chown=user:user /project/cmd/${SERVICE}/app /cmd/app
COPY --from=builder --chown=user:user \
    /project/${APP_PATH}/config/config.ini \
    /cmd/${APP_PATH}/config/config.ini
COPY --from=builder --chown=user:user \
    /project/${APP_PATH}/log/log.txt \
    /cmd/${APP_PATH}/log/log.txt
USER user
EXPOSE ${PORT}
ENTRYPOINT ["/cmd/app"]