FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY db.go main.go handlers.go hub.go models.go ./
RUN go build -o api .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 3000
CMD ["./api"]