# Makefile for In the Light of Love

# Variables - These can be overridden by environment variables
REGISTRY ?= registry.digitalocean.com/praslar
IMAGE_NAME ?= in-the-light-of-love
TAG ?= latest
IMAGE_TAG = $(REGISTRY)/$(IMAGE_NAME):$(TAG)

.PHONY: all build push run down

all: build

# Build the Docker image for multi-platform
build:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_TAG) .

# Push the Docker image to the registry
push:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_TAG) --push .

# Run the application using Docker Compose
run:
	docker-compose pull
	docker-compose up -d

# Stop and remove the application containers
down:
	docker-compose down
