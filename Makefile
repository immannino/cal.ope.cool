.PHONY: openapi
openapi:
	@./scripts/openapi.sh nhl

.PHONY: sqlc
sqlc:
	sqlc generate