# builder stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git bash
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gateway .

# final small image
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/gateway .
EXPOSE 8000
CMD ["./gateway"]
