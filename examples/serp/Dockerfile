FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go mod init serp-scraper && go build -o serp-scraper .

FROM alpine:latest

WORKDIR /output
COPY --from=builder /app/serp-scraper /usr/local/bin/serp-scraper

ENTRYPOINT ["serp-scraper"]