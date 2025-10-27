FROM golang:1.24.9-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY ./go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o user-activity-tracking-api ./cmd

FROM alpine:3.16

COPY --from=builder /app/user-activity-tracking-api /user-activity-tracking-api

ENTRYPOINT ["/user-activity-tracking-api"]