#!/bin/bash

# 1. Generate Swagger 2.0 specs
go generate ./...

# 2. Delete unused docs.go to prevent IDE import errors
rm -f docs/docs.go

# 3. Convert Swagger 2.0 to OpenAPI 3.0 (YAML)
npx -y swagger2openapi --yaml --outfile docs/openapi.yaml docs/swagger.json

# 4. Convert OpenAPI 3.0 (YAML) to OpenAPI 3.1 (YAML)
npx -y openapi-format docs/openapi.yaml -o docs/openapi.yaml --convertTo "3.1"

# 5. Convert OpenAPI 3.1 (YAML) to OpenAPI 3.1 (JSON)
npx -y openapi-format docs/openapi.yaml -o docs/openapi.json --convertTo "3.1"

# 6. Build the Go project
go build

