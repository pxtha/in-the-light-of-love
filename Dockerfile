# Stage 1: Build the Go backend
FROM golang:1.24.4-alpine AS go-builder

# Enable CGO and install build tools for the sqlite driver
ENV CGO_ENABLED=1 GOOS=linux
RUN apk add --no-cache build-base

ARG TARGETARCH
WORKDIR /app

COPY go.mod ./
# go mod download will create go.sum if it doesn't exist
RUN go mod download

COPY . .
RUN GOARCH=${TARGETARCH} go build -o /in-the-light-of-love .

# Stage 2: Final image
FROM alpine:latest

# Add compatibility libraries for CGO binaries
RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY --from=go-builder /in-the-light-of-love .
COPY --from=go-builder /app/static ./static
COPY --from=go-builder /app/templates ./templates
COPY --from=go-builder /app/uploads ./uploads
COPY gallery.db .

EXPOSE 8080
CMD ["/app/in-the-light-of-love"]
