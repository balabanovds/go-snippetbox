progname=snippetbox

include .env

dev:
	$(GO.EXE) run ./cmd/web -db_user=$(DB_USER) -db_password=$(DB_PASSWORD) \
		-db_host=$(DB_HOST) -db_port=$(DB_PORT)
