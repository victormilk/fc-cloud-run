FROM golang:1.23.3 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/cloud-run/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/.env .
COPY --from=builder /app/main .
CMD ["./main"]