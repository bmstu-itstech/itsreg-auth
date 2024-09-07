.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh auth internal/ports/httpport httpport
