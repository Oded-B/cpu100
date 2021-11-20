##
## Build stage
##

# Use image with golang last version as builder.
FROM golang:1.17-bullseye AS build

# Make project root folder as current dir.
WORKDIR $GOPATH/src/github.com/schwarzlichtbezirk/cpu100/
# Copy only go.mod and go.sum to prevent downloads all dependencies again on any code changes.
COPY go.* .
# Download all dependencies pointed at go.mod file.
RUN go mod download
# Copy all files and subfolders in current state as is.
COPY . .

# Build service and put executable.
RUN go build -o /go/bin/cpu100 .

##
## Deploy stage
##

# Thin deploy image.
FROM gcr.io/distroless/base-debian11

# Copy compiled executables to new image destination.
COPY --from=build /go/bin/cpu100 /go/bin/cpu100

# Run application with full path representation.
# Without shell to get signal for graceful shutdown.
ENTRYPOINT ["/go/bin/cpu100"]
