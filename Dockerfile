# Stage 1: Build the Go backend
FROM golang:1.24.4-alpine AS go-builder

ENV CGO_ENABLED=0 GOOS=linux
ARG TARGETARCH
WORKDIR /app

COPY go.mod ./
# go mod download will create go.sum if it doesn't exist
RUN go mod download

COPY . .
RUN GOARCH=${TARGETARCH} go build -o /in-the-light-of-love .

# Stage 2: Final image
FROM alpine:latest

WORKDIR /app
COPY --from=go-builder /in-the-light-of-love .
COPY --from=go-builder /app/static ./static
COPY --from=go-builder /app/templates ./templates
COPY --from=go-builder /app/uploads ./uploads

EXPOSE 8080
CMD ["/app/in-the-light-of-love"]
