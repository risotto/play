  
# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13 as builder

# Add Maintainer Info
LABEL maintainer="James Jarvis <git@jamesjarvis.io>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o api cmd/play/main.go


######## Start a new stage from scratch #######
FROM raphaelvigee/risotto:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api .
RUN chmod +x ./api

EXPOSE 4000

ENV PATH="/:${PATH}"

ENTRYPOINT [ ]

# Command to run the executable
CMD ["./api"] 