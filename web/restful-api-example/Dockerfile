# syntax=docker/dockerfile:1

# Use official Go image as the base image. It contains
# environment for building Go programs.
FROM golang:1.17-alpine AS build

# Create a directory inside the image. Docker will use this directory as the
# default desitnation for all subsequent commands (you can use relative paths).
WORKDIR /app

# Download all modules necessary to compile the app, i.e. download dependencies.
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Compile (build the source code of) the app.
COPY *.go ./
RUN go build -o /restful-api-example

# Create a single layer image.
FROM alpine:latest
COPY --from=build /restful-api-example /restful-api-example

# Doesn't actually publish the port but functions as documentation. To actually
# publish the port use `-p` on `docker run` or `-P` to publish all exposed ports
# and map them to high-order ports.
EXPOSE 8080

# Tell Docker to execute this command on a "docker run". I.e. run
# the app when the container is started.
CMD [ "/restful-api-example" ]
