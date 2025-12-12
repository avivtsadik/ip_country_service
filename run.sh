#!/bin/bash

# Docker run script for IP Country Service

set -e

IMAGE_NAME="ip-country-service"
TAG="latest"
CONTAINER_NAME="ip-country-service"
HOST_PORT="8080"
CONTAINER_PORT="8080"

# Stop and remove existing container if running
if docker ps -a --format 'table {{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "Stopping and removing existing container: ${CONTAINER_NAME}"
    docker stop "${CONTAINER_NAME}" >/dev/null 2>&1 || true
    docker rm "${CONTAINER_NAME}" >/dev/null 2>&1 || true
fi

echo "Starting Docker container: ${CONTAINER_NAME}"
echo "Image: ${IMAGE_NAME}:${TAG}"
echo "Port mapping: ${HOST_PORT}:${CONTAINER_PORT}"

# Run the container
docker run -d \
    --name "${CONTAINER_NAME}" \
    -p "${HOST_PORT}:${CONTAINER_PORT}" \
    "${IMAGE_NAME}:${TAG}"

echo "✅ Container started successfully!"
echo ""
echo "Service is running at: http://localhost:${HOST_PORT}"
echo ""
echo "Test the service:"
echo "  curl \"http://localhost:${HOST_PORT}/health\""
echo "  curl \"http://localhost:${HOST_PORT}/v1/find-country?ip=8.8.8.8\""
echo ""
echo "View container logs:"
echo "  docker logs ${CONTAINER_NAME}"
echo ""
echo "Stop the container:"
echo "  docker stop ${CONTAINER_NAME}"