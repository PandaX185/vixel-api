#!/bin/bash

# Script to test all API endpoints
# Starts the server, creates user, uploads image, tests all endpoints, generates markdown docs
# Requires: curl, jq, imagemagick

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting API endpoint tests...${NC}"

# Check dependencies
echo "Checking dependencies..."
command -v curl >/dev/null 2>&1 || { echo -e "${RED}curl is required but not installed.${NC}"; exit 1; }
command -v jq >/dev/null 2>&1 || { echo -e "${RED}jq is required but not installed.${NC}"; exit 1; }
command -v convert >/dev/null 2>&1 || { echo -e "${RED}ImageMagick (convert) is required but not installed.${NC}"; exit 1; }

# Check if PostgreSQL is accessible
echo "Checking PostgreSQL connection..."
export PGPASSWORD="root"
if ! psql -h localhost -p 5432 -U postgres -d postgres -c "SELECT 1;" >/dev/null 2>&1; then
    echo -e "${RED}PostgreSQL is not accessible. Please ensure PostgreSQL is running and credentials are correct.${NC}"
    echo -e "${YELLOW}You may need to run: createdb vixel_test${NC}"
    exit 1
fi

# Set environment variables for testing
export DB_URL="postgres://postgres:root@localhost:5432/vixel_test?sslmode=disable"
export JWT_SECRET="test_secret_key_for_integration_testing"
export MINIO_ENDPOINT="localhost:9000"
export MINIO_ACCESS_KEY="minioadmin"
export MINIO_SECRET_KEY="minioadmin"
export MINIO_BUCKET_NAME="vixel"
export MINIO_USE_SSL="false"
export MINIORegion="us-east-1"

echo -e "${YELLOW}Prerequisites:${NC}"
echo -e "${YELLOW}- PostgreSQL running on localhost:5432 with user 'postgres' and password 'root'${NC}"
echo -e "${YELLOW}- MinIO running on localhost:9000 with default credentials${NC}"
echo -e "${YELLOW}Create database: createdb vixel_test${NC}"

# Reset database
echo "Resetting test database..."
export PGPASSWORD="root"
# Terminate active connections to the database
psql -h localhost -p 5432 -U postgres -d postgres -c "
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE datname = 'vixel_test' AND pid <> pg_backend_pid();" >/dev/null 2>&1
psql -h localhost -p 5432 -U postgres -d postgres -c "DROP DATABASE IF EXISTS vixel_test;" >/dev/null 2>&1
psql -h localhost -p 5432 -U postgres -d postgres -c "CREATE DATABASE vixel_test;" >/dev/null 2>&1

# Start compose file 
docker-compose down
docker-compose up -d

# Start the server
echo "Starting server..."
# Kill any existing server on port 8080
pkill -f "go run app/main.go" || true
sleep 2

go run app/main.go &
SERVER_PID=$!
sleep 5  # Wait for server to start

# Check if server is running
echo "Checking if server is running..."
if ! curl -s --max-time 5 http://localhost:8080/api/v1/users > /dev/null; then
    echo -e "${RED}Server failed to start or API is not responding on localhost:8080/api/v1${NC}"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi
echo -e "${GREEN}Server is running.${NC}"

# Create a test image
echo "Creating test image..."
convert -size 100x100 xc:red test.jpg 2>/dev/null || {
    echo -e "${RED}ImageMagick not available. Please install it or create test.jpg manually.${NC}"
    kill $SERVER_PID
    exit 1
}

# Register user
echo "Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","email":"test@example.com","password":"password123"}')

echo "# Vixel API Documentation" > api_documentation.md
echo "" >> api_documentation.md
echo "## Register User" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/users" >> api_documentation.md
echo 'Body: {"username":"testuser","email":"test@example.com","password":"password123"}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$REGISTER_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Login
echo "Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123"}')

echo "## Login" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/users/login" >> api_documentation.md
echo 'Body: {"email":"test@example.com","password":"password123"}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$LOGIN_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}Failed to get token${NC}"
    exit 1
fi

echo "Token obtained: ${TOKEN:0:20}..."

# Function to upload test image and get ID
upload_test_image() {
    local desc="$1"
    echo "Uploading test image for $desc..." >&2
    local response=$(curl -s -X POST http://localhost:8080/api/v1/images \
        -H "Authorization: Bearer $TOKEN" \
        -F "file=@test.jpg" \
        -F "alt_text=$desc")
    local image_id=$(echo $response | jq -r '.data.id')
    echo $image_id
}

# Upload image
echo "Uploading image..."
UPLOAD_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images \
    -H "Authorization: Bearer $TOKEN" \
    -F "file=@test.jpg" \
    -F "alt_text=Test Image")

echo "## Upload Image" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}" >> api_documentation.md
echo "Form: file=@test.jpg, alt_text=Test Image" >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$UPLOAD_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Get image test
echo "Testing get image..."
IMAGE_ID_GET=$(upload_test_image "Get Test")
GET_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/images/$IMAGE_ID_GET \
    -H "Authorization: Bearer $TOKEN")

echo "## Get Image" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "GET /api/v1/images/{id}" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}" >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$GET_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# List user images
echo "Testing list user images..."
LIST_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/users/1/images \
    -H "Authorization: Bearer $TOKEN")

echo "## List User Images" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "GET /api/v1/users/{user_id}/images" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}" >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$LIST_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Transformations
echo "Testing transformations..."
echo "" >> api_documentation.md
echo "# Image Transformations" >> api_documentation.md
echo "" >> api_documentation.md

# Resize
echo "Testing resize transformation..."
IMAGE_ID_RESIZE=$(upload_test_image "Resize Test")
RESIZE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_RESIZE/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"resize":{"width":50,"height":50}}')

echo "## Transform Image - Resize" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"resize":{"width":50,"height":50}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$RESIZE_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Crop
echo "Testing crop transformation..."
IMAGE_ID_CROP=$(upload_test_image "Crop Test")
CROP_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_CROP/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"crop":{"x":10,"y":10,"width":50,"height":50}}')

echo "## Transform Image - Crop" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"crop":{"x":10,"y":10,"width":50,"height":50}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$CROP_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Rotate
echo "Testing rotate transformation..."
IMAGE_ID_ROTATE=$(upload_test_image "Rotate Test")
ROTATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_ROTATE/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"rotate":{"angle":90}}')

echo "## Transform Image - Rotate" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"rotate":{"angle":90}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$ROTATE_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Flip Horizontal
echo "Testing flip horizontal transformation..."
IMAGE_ID_FLIP_H=$(upload_test_image "Flip Horizontal Test")
FLIP_H_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_FLIP_H/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"flip":{"direction":"horizontal"}}')

echo "## Transform Image - Flip Horizontal" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"flip":{"direction":"horizontal"}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$FLIP_H_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Flip Vertical
echo "Testing flip vertical transformation..."
IMAGE_ID_FLIP_V=$(upload_test_image "Flip Vertical Test")
FLIP_V_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_FLIP_V/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"flip":{"direction":"vertical"}}')

echo "## Transform Image - Flip Vertical" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"flip":{"direction":"vertical"}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$FLIP_V_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md

# Format Conversion
echo "Testing format conversion transformation..."
IMAGE_ID_FORMAT=$(upload_test_image "Format Conversion Test")
FORMAT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_FORMAT/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"format_conversion":{"format":"png"}}')

echo "## Transform Image - Format Conversion" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"format_conversion":{"format":"png"}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$FORMAT_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md
echo "" >> api_documentation.md

# Filter
echo "Testing filter transformation..."
IMAGE_ID_FILTER=$(upload_test_image "Filter Test")
FILTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_FILTER/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"filter":{"saturation":50,"brightness":10,"contrast":20}}')

echo "## Transform Image - Filter" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"filter":{"saturation":50,"brightness":10,"contrast":20}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$FILTER_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md
echo "" >> api_documentation.md

# Watermark
echo "Testing watermark transformation..."
IMAGE_ID_WATERMARK=$(upload_test_image "Watermark Test")
WATERMARK_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/images/$IMAGE_ID_WATERMARK/transform \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"watermark":{"text":"Test","position":{"x":10,"y":10},"opacity":50}}')

echo "## Transform Image - Watermark" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "POST /api/v1/images/{id}/transform" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}, Content-Type: application/json" >> api_documentation.md
echo 'Body: {"watermark":{"text":"Test","position":{"x":10,"y":10},"opacity":50}}' >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo "$WATERMARK_RESPONSE" >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md
echo "" >> api_documentation.md

# Delete test images
echo "Cleaning up test images..."
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_GET -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_RESIZE -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_CROP -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_ROTATE -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_FLIP_H -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_FLIP_V -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_FORMAT -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_FILTER -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1
curl -s -X DELETE http://localhost:8080/api/v1/images/$IMAGE_ID_WATERMARK -H "Authorization: Bearer $TOKEN" >/dev/null 2>&1

echo "## Delete Image" >> api_documentation.md
echo "### Request" >> api_documentation.md
echo "DELETE /api/v1/images/{id}" >> api_documentation.md
echo "Headers: Authorization: Bearer {token}" >> api_documentation.md
echo "### Response" >> api_documentation.md
echo "\`\`\`json" >> api_documentation.md
echo '{"data":{"message":"image deleted"},"status":"success","timestamp":"2026-02-05T07:13:24.180781809+02:00"}' >> api_documentation.md
echo "\`\`\`" >> api_documentation.md
echo "" >> api_documentation.md
echo "" >> api_documentation.md

# Cleanup
rm -f test.jpg
rm -f register.md login.md upload.md get_image.md list_images.md
rm -f transform_*.md delete_image.md

echo -e "${GREEN}All tests completed. API documentation generated in api_documentation.md${NC}"

# Stop server
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true