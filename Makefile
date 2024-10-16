include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -limiter-enabled=false -cors-trusted-origins='http://localhost:9000'

## db/start: compose up and connect to database using psql
.PHONY: db/start
db/start:
	docker compose rm -f && docker compose pull && docker compose up -d && clear && printf 'All done!\n\nEnter `psql -h localhost -p 5432 -d householdingindex -U admin` to start\n' && docker exec -it postgres_container sh

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN} up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${DB_DSN} down

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
# COULD INCLUDE STATIC CHECK AND TESTS

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

# ==================================================================================== #
# POPULATING DATABASE WITH SAMPLE DATA
# ==================================================================================== #

## pop/createadmin: create an admin user
.PHONY: pop/createadmin
pop/createadmin:
	"./shell scripts/populate database/populate_database_createadmin.sh"

## pop/getadmintoken: get token from admin user
.PHONY: pop/getadmintoken
pop/getadmintoken:
	"./shell scripts/populate database/populate_database_getadmintoken.sh"

## pop/itemtypes: create some itemtypes
.PHONY: pop/itemtypes
pop/itemtypes:
	"./shell scripts/populate database/populate_database_itemtypes.sh"

## pop/measurements: create some measurements
.PHONY: pop/measurements
pop/measurements:
	"./shell scripts/populate database/populate_database_measurements.sh"

## pop/knownitems: create some known items
.PHONY: pop/knownitems
pop/knownitems:
	"./shell scripts/populate database/populate_database_knownitems.sh"

## pop/availableitems: create some available items
.PHONY: pop/availableitems
pop/availableitems:
	"./shell scripts/populate database/populate_database_availableitems.sh"

## pop/ingredients: create some ingredients
.PHONY: pop/ingredients
pop/ingredients:
	"./shell scripts/populate database/populate_database_ingredients.sh"

## pop/recipies: create some recipies
.PHONY: pop/recipies
pop/recipies:
	"./shell scripts/populate database/populate_database_recipies.sh"

## pop/recipeingredients: create some recipeingredients
.PHONY: pop/recipeingredients
pop/recipeingredients:
	"./shell scripts/populate database/populate_database_recipeingredients.sh"