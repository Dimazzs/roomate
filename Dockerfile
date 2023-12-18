# Build stage
FROM golang:1.21.5-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main app.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY start.sh .
COPY wait-for.sh .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]