  
# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13 as builder

# Add Maintainer Info
LABEL maintainer="James Jarvis <git@jamesjarvis.io>"
ENV CGO_ENABLED 0

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o api cmd/play/main.go

FROM builder as tester

COPY --from=raphaelvigee/risotto:latest /rst /usr/bin/rst

CMD [ "go", "test", "/app/..." ]

FROM debian as runner

COPY --from=raphaelvigee/risotto:latest /rst /usr/bin/rst
COPY --from=builder /app/api .

CMD ["./api"] 
