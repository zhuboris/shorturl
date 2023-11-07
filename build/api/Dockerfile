FROM golang:1.21.3-alpine3.18 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o shorturl_api ./cmd/

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /app/shorturl_api .

EXPOSE 3000
EXPOSE 50051

CMD ["./shorturl_api"]