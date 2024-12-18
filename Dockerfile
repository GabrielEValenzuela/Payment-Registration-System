# Stage 1: Build the Go application
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code
COPY ./src ./src

# Build the Go executable
RUN go build -o main ./src/cmd/main.go

# Stage 2: Runtime environment
FROM ubuntu:24.04 

# Set the working directory inside the container
WORKDIR /app

# Install netcat (nc) for the wait-for.sh script
RUN apt-get update && apt-get install -y netcat-openbsd && apt-get clean

# Copy the compiled executable from the build stage
COPY --from=builder /app/main .

# Copy the configuration file
COPY ./config_docker.yml ./config.yml

# Copy the wait-for.sh script
COPY wait-for.sh /wait-for.sh
RUN chmod +x /wait-for.sh

# Copy sql demo data
COPY ./src/internal/storage/relational/insert.sql ./src/internal/storage/relational/insert.sql

# Expose the application port
EXPOSE 8080

# Use wait-for.sh to ensure dependencies are ready before starting the application
CMD ["/wait-for.sh", "./main"]
