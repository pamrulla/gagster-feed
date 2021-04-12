FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CG0_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Run test
RUN go test ./... 

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Export necessary port
# EXPOSE 3000

# Command to run when starting the container
CMD ["/dist/main"]

# Build a small image
# FROM scratch

# COPY --from=builder /dist/ /

# Command to run the executable
# ENTRYPOINT ["/main"]
