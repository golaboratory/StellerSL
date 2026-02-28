#!/bin/bash

# Configuration
API_URL="http://localhost:8888/openapi.json"
OUTPUT_DIR="./frontend/src/api"
OPENAPI_JSON="./frontend/openapi.json"

echo "Fetching OpenAPI spec from $API_URL..."
curl -s $API_URL -o $OPENAPI_JSON

if [ ! -s "$OPENAPI_JSON" ]; then
    echo "Error: Failed to fetch openapi.json from backend. Is it running?"
    exit 1
fi

echo "Generating TypeScript API client using Docker..."
mkdir -p $OUTPUT_DIR

# Use Docker to run openapi-generator-cli
# We map the current directory to /local in the container
# Using v7.18.0 which is available locally
docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli:v7.18.0 generate \
  -i /local/frontend/openapi.json \
  -g typescript-axios \
  -o /local/frontend/src/api \
  --additional-properties=useSingleRequestParameter=true,supportsES6=true

echo "Done! API objects generated in $OUTPUT_DIR"
rm frontend/openapi.json
