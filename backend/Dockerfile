FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ARG GITHUB_REPO=kelvinwambua/nimbledb
ARG GITHUB_BRANCH=master

RUN git clone https://github.com/${GITHUB_REPO}.git . && \
    git checkout ${GITHUB_BRANCH} || git checkout master || git checkout main

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o nimbledb-server main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/nimbledb-server .

RUN mkdir -p /data

EXPOSE 7000

CMD ["./nimbledb-server", "-mode", "server", "-addr", ":7000", "-data", "/data"]
