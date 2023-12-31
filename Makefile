postgres:
	docker run --name dropbyte-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.3-alpine
createdb:
	docker exec -it dropbyte-db createdb --username=root --owner=root dropbyte
migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/dropbyte?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/dropbyte?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=dropbyte \
  proto/*.proto
evans:
	evans --host localhost --port 9090 -r repl

.PHONY: sqlc postgres createdb test server proto evans
