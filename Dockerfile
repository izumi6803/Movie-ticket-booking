FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY cinema-backend/go.mod ./
COPY cinema-backend/go.sum* ./
RUN go mod download

COPY cinema-backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 3001
ENV PORT=3001
CMD ["./main"]
