#!/bin/bash

# Exit on error
set -e

# Colors for output
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Glimpse Docker images...${NC}"

# Build agent image
echo -e "\n${GREEN}Building agent image...${NC}"
docker build -t glimpse-agent:latest -f docker/Dockerfile.agent .

# Build server image
echo -e "\n${GREEN}Building server image...${NC}"
docker build -t glimpse-server:latest -f docker/Dockerfile.server .

echo -e "\n${GREEN}Build complete!${NC}"
echo -e "Images created:"
echo -e "- glimpse-agent:latest"
echo -e "- glimpse-server:latest" 