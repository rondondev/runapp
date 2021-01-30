MIGRATION_CONN_STRING="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
MIGRATION_CONN_STRING_TEST="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:6433/${POSTGRES_DB}?sslmode=disable"

postgres:
	docker run --name runapp-postgres -p ${POSTGRES_PORT}:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:13-alpine

createdb:
	docker exec -it runapp-postgres createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}

dropdb:
	docker exec -it runapp-postgres dropdb ${POSTGRES_DB}

migrateup:
	migrate -path db/migration -database ${MIGRATION_CONN_STRING} -verbose up

migratedown:
	migrate -path db/migration -database ${MIGRATION_CONN_STRING} -verbose down

dbmigration:
	migrate create -ext sql -dir db/migration -seq $(m)

createtestdb:
	@docker stop runapp-pg-test || true
	@docker rm runapp-pg-test || true
	@docker run --name runapp-pg-test -p 6433:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:13-alpine
	./wait-for-pg.sh localhost 6433
	@docker exec -it runapp-pg-test createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}
	@migrate -path db/migration -database ${MIGRATION_CONN_STRING_TEST} -verbose up

 sqlc:
	sqlc generate

test:
	go test -v -cover ./...

run:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test serve