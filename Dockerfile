#Stage 1: Download deps & build
FROM golang:alpine as builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build
#Stage 2: run
FROM golang:alpine
WORKDIR /app
COPY --from=builder /go/src/app/loki-index-regen app
CMD ["./app"]