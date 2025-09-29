include .env
export

MIGRATION_DIRS = src/db/migrations

# Create a new migration (make migrate-create NAME=profiles)
migrate-create: 
	migrate create -ext sql -dir $(MIGRATION_DIRS) -seq $(NAME)

# Run all pending migration (migrate-up)
migrate-up:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" up

# Rollback the last migration (migrate-down)
migrate-down:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down 1

# Rollback n migration (migrate-down-n)
migrate-down-n:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" down $(N)

# Forcr migration version (use with caution - make migrate-force VERSION=1)
migrate-force:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" force $(VERSION)

# Drop all (include schema migration - )
migrate-drop:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" drop

# Apply specific migration version (make migrate-goto VERSION=1)
migrate-goto:
	migrate -path $(MIGRATION_DIRS) -database "$(CONN_STRING)" goto $(VERSION)

debug:
	@echo $(CONN_STRING)

.PHONY: migrate-create migrate-up debug migrate-force migrate-drop migrate-goto migrate-down-n