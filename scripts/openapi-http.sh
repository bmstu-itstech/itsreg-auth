#!/bin/bash
set -e

readonly service="$1"
readonly output_dir="$2"
readonly package="$3"

oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yaml"
oapi-codegen -generate chi-server -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yaml"
mkdir -p "api/openapi/clients/$service"
oapi-codegen -generate types -o "api/openapi/clients/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yaml"
oapi-codegen -generate client -o "api/openapi/clients/$service/openapi_client_gen.go" -package "$service" "api/openapi/$service.yaml"
