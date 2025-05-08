SQLC := /Users/vitaliiiavurek/go/bin/sqlc
MOCKGEN = /Users/vitaliiiavurek/go/bin/mockgen

.PHONY: postgres
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

.PHONY: createdb
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

.PHONY: dropdb
dropdb:
	docker exec -it postgres dropdb simple_bank

.PHONY: migrate
migrate:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" up

.PHONY: rollback
rollback:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" down

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