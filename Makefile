postgres:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.3-alpine
createdb:
	docker exec -it postgres15 createdb --username=root --owner=root dropbyte
migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/dropbyte?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/dropbyte?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go

.PHONY: sqlc postgres createdb test server
