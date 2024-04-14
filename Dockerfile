FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o "config-creator" .

FROM alpine:latest

# ca-certificates are required for TLS connections within static binary applications
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache nginx && \
    mkdir -p /run/nginx

RUN rm /etc/nginx/nginx.conf

WORKDIR /root/
COPY --from=builder /app/config-creator .
ENTRYPOINT /root/config-creator
