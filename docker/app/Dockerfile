# Build stage
FROM golang:1.20.7-alpine3.18 AS build

# Add Maintainer Info
LABEL maintainer="Thais Carvalho <email>"

WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed 
RUN go mod download 

# Copy source files
COPY pkg ./pkg
COPY cmd ./cmd
COPY utils ./utils

# Change to the directory containing main.go
WORKDIR /app/cmd

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM scratch

# Import the user and group files from the builder.
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

# Copy binary from build to main stage
COPY --from=build /app/cmd/main .

# Use an unprivileged user.
USER nobody:nobody

# Export necessary port
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]
