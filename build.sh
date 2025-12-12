#!/bin/bash

# Docker build script for IP Country Service

set -e

IMAGE_NAME="ip-country-service"
TAG="latest"

echo "Building Docker image: ${IMAGE_NAME}:${TAG}"

# Build the Docker image
docker build -t "${IMAGE_NAME}:${TAG}" .

echo "✅ Docker image built successfully!"
echo "Image: ${IMAGE_NAME}:${TAG}"
echo ""
echo "To run the container:"
echo "  ./run.sh"
echo ""
echo "To run manually:"
echo "  docker run -p 8080:8080 ${IMAGE_NAME}:${TAG}"