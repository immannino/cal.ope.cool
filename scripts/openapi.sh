#!/bin/bash
set -e

readonly service="$1"
# readonly output_dir="$2"
# readonly package="$3"

# oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/$service-api.yml"
# oapi-codegen -generate gorilla,spec -o "$output_dir/openapi_api.gen.go" -package "$package" "api/$service-api.yml"
oapi-codegen -generate types -o "pkg/$service/openapi_types.gen.go" -package "$service" "api/$service.yaml"
oapi-codegen -generate client -o "pkg/$service/openapi_client_gen.go" -package "$service" "api/$service.yaml"
