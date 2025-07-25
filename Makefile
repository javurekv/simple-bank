DB_URL = postgres://root:secret@localhost:5432/simple_bank?sslmode=disable
SQLC := /Users/vitaliiiavurek/go/bin/sqlc
MOCKGEN = /Users/vitaliiiavurek/go/bin/mockgen

.PHONY: postgres
postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

.PHONY: createdb
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

.PHONY: dropdb
dropdb:
	docker exec -it postgres dropdb simple_bank

.PHONY: migrate
migrate:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "$(DB_URL)" up

.PHONY: migrate1
migrate1:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "$(DB_URL)" up 1

.PHONY: rollback
rollback:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "$(DB_URL)" down

.PHONY: rollback1
rollback1:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "$(DB_URL)" down 1

.PHONY: dbdocs
dbdocs:
	dbdocs build doc/db.dbml

.PHONY: dbschema
dbschema:
	dbml2sql --postgres -o doc/schema.sql doc/bd.dbml

.PHONY: sqlc
sqlc:
	$(SQLC) generate

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: server
server:
	go run main.go

.PHONY: mock
mock:
	$(MOCKGEN) -package mockdb -destination=db/mock/mock_store.go simple_bank.sqlc.dev/app/db/sqlc Store

.PHONY: proto

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
		proto/*.proto
	statik -src=./doc/swagger -dest=./doc

.PHONY: evans
evans:
	evans --host localhost --port 9090 -r relp
