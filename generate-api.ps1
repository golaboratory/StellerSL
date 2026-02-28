# Configuration
$API_URL = "http://localhost:8888/openapi.json"
$OUTPUT_DIR = "./frontend/src/api"
$OPENAPI_JSON = "./frontend/openapi.json"

Write-Host "Fetching OpenAPI spec from $API_URL..." -ForegroundColor Cyan
try {
    Invoke-WebRequest -Uri $API_URL -OutFile $OPENAPI_JSON
} catch {
    Write-Host "Error: Failed to fetch openapi.json from backend. Is it running?" -ForegroundColor Red
    exit 1
}

if (-Not (Test-Path $OPENAPI_JSON)) {
    Write-Host "Error: openapi.json not found." -ForegroundColor Red
    exit 1
}

Write-Host "Generating TypeScript API client..." -ForegroundColor Cyan
if (-Not (Test-Path $OUTPUT_DIR)) {
    New-Item -ItemType Directory -Force -Path $OUTPUT_DIR
}

Set-Location -Path "./frontend"
pnpm openapi-generator-cli generate `
  -i "../frontend/openapi.json" `
  -g typescript-axios `
  -o "./src/api" `
  --additional-properties=useSingleRequestParameter=true,supportsES6=true

Set-Location -Path ".."
Remove-Item -Path $OPENAPI_JSON

Write-Host "Done! API objects generated in $OUTPUT_DIR" -ForegroundColor Green
