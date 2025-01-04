#!/bin/bash

# this is for making the docker image work on the prod environemnt. The server is hosted on digitalocean with low memory, therefore I need to add swap space

# Set strict mode
set -euo pipefail

# Function to log messages
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

# Function to clean up swap
cleanup_swap() {
    log "Cleaning up swap space..."
    swapoff /swapfile || log "Failed to turn off swap"
    rm /swapfile || log "Failed to remove swapfile"
}

# Trap to ensure swap cleanup on script exit
trap cleanup_swap EXIT

# Get the name of the user who invoked sudo
SUDO_USER="${SUDO_USER:-$USER}"

# Function to run git commands as the original user
git_pull_as_user() {
    sudo -u "$SUDO_USER" git -C "$1" pull
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    log "Please run with sudo"
    exit 1
fi

# Add swap space
log "Adding swap space..."
fallocate -l 2G /swapfile || { log "Failed to allocate swap"; exit 1; }
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile

# Check if .env.prod file exists
if [ ! -f .env.prod ]; then
    log "Error: .env.prod file not found"
    exit 1
fi

# Pull the latest changes from main in the current directory
log "Pulling latest changes from main in the current directory..."
git_pull_as_user "$(pwd)" || { log "Failed to pull latest changes in the current directory"; exit 1; }

# Pull the latest changes from main in the ../indicum-frontend directory
log "Pulling latest changes from main in the ../indicum-frontend directory..."
git_pull_as_user "../indicum-frontend" || { log "Failed to pull latest changes in ../indicum-frontend"; exit 1; }

# Pull the latest changes from main in the ../indicum-server directory
log "Pulling latest changes from main in the ../indicum-server directory..."
git_pull_as_user "../indicum-server" || { log "Failed to pull latest changes in ../indicum-server"; exit 1; }

# Build Docker images
log "Building Docker images..."
docker compose --env-file .env.prod build || { log "Docker build failed"; exit 1; }

# Start Docker containers
log "Starting Docker containers..."
docker compose --env-file .env.prod up -d || { log "Docker up failed"; exit 1; }

log "Script completed successfully"