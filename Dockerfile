# Stage 1: Build the Go application
FROM golang:latest AS build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Install curl
RUN apk add --no-cache curl

# Copy the .env file
COPY .env .env

# Copy the built Go binary from the build stage
COPY --from=build /app/main .

EXPOSE 8081

CMD ["./main"]


