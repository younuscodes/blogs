# Use an official Golang runtime as a parent image
FROM golang:1.16.5-alpine3.13 AS build-env

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o app

# Use a lightweight Alpine image as a base
FROM alpine:3.13

# Copy the executable file from the previous build step
COPY --from=build-env /app/app /

# Expose port 8123 to the outside world
EXPOSE 8123

# Set the default command to run the application
CMD ["/app"]
