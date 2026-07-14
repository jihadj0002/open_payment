FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/mock-processor .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/mock-processor .

EXPOSE 9000

HEALTHCHECK --interval=10s --timeout=3s --start-period=2s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:9000/health || exit 1

ENTRYPOINT ["./mock-processor"]
