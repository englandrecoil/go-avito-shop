FROM golang:1.23.1-alpine AS builder

WORKDIR /build
COPY . . 
RUN go mod download
RUN go build -o ./shop-service

FROM gcr.io/distroless/base-debian12
WORKDIR /build
COPY --from=builder /build/shop-service ./shop-service
COPY .env .env
CMD ["/build/shop-service"]

