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
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "postgres://root:Kkngb2MMWjynXJNugtAZW@simple-bank.c9mimg8mqt7h.eu-north-1.rds.amazonaws.com:5432/simple_bank" up

.PHONY: migrate1
migrate1:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "postgres://root:Kkngb2MMWjynXJNugtAZW@simple-bank.c9mimg8mqt7h.eu-north-1.rds.amazonaws.com:5432/simple_bank" up 1

.PHONY: rollback
rollback:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" down

.PHONY: rollback1
rollback1:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cli@v4.17.0 -path ./db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" down 1

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