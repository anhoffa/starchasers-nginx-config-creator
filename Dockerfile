FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o "config-creator" .

FROM alpine:latest

# ca-certificates are required for TLS connections within static binary applications
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache nginx nginx-mod-stream curl && \
    mkdir -p /run/nginx

RUN rm /etc/nginx/nginx.conf

WORKDIR /root/
COPY --from=builder /app/config-creator .
COPY nginxTemplate.tmpl /root/nginxTemplate.tmpl

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 CMD curl --fail http://localhost:7070/healthz || exit 1
ENTRYPOINT /root/config-creator
