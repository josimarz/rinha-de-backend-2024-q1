FROM golang:alpine AS builder

WORKDIR /var/app

COPY . .

RUN go build cmd/server/main.go

FROM scratch

ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG REDIS_ADDR
ARG REDIS_PORT
ARG REDIS_PASSWORD
ARG REDIS_DB

WORKDIR /var/app

COPY --from=builder /var/app/main .

EXPOSE 8080

ENTRYPOINT [ "./main" ]