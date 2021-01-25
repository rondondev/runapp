MIGRATION_CONN_STRING="postgresql://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable"

postgres:
	docker run --name runapp-postgres -p ${DB_PORT}:5432 -e POSTGRES_USER=${DB_USERNAME} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:13-alpine

createdb:
	docker exec -it runapp-postgres createdb --username=${DB_USERNAME} --owner=${DB_USERNAME} ${DB_NAME}

dropdb:
	docker exec -it runapp-postgres dropdb ${DB_NAME}

migrateup:
	migrate -path db/migration -database ${MIGRATION_CONN_STRING} -verbose up

migratedown:
	migrate -path db/migration -database ${MIGRATION_CONN_STRING} -verbose down

dbmigration:
	migrate create -ext sql -dir db/migration -seq $(m)

 sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test