MIGRATIONS_DIR := db/migrations

start_server:
	go run ./cmd/web -port=:${PORT} -dsn=${POSTGRES_DSN}

create_migrations:
	migrate create -ext sql -dir ${MIGRATIONS_DIR} -seq create_snippets_table

run_migrations:
	migrate -database ${POSTGRES_DSN} -path ${MIGRATIONS_DIR} up

roll_migrations:
	migrate -database ${POSTGRES_DSN} -path ${MIGRATIONS_DIR} down

force_migrate:
	migrate -database ${POSTGRES_DSN} -path ${MIGRATIONS_DIR} force ${VERSION} 