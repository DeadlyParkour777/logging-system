FROM golang:1.24-alpine AS builder
ARG SERVICE_PATH

WORKDIR /build

COPY go.work go.work.sum ./
COPY services/ services/
COPY pkg/ pkg/

RUN go work vendor

COPY . .

RUN go build -o /app/server ./services/${SERVICE_PATH}/cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]