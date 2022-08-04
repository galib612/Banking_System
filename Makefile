mariadb:
	docker run --detach --name mariadbauth -p 3306:3306 -e MARIADB_USER=galib -e MARIADB_PASSWORD=secret -e MARIADB_ROOT_PASSWORD=secret mariadb:latest

mariadbauthstart:
	docker start mariadbauth

mariadbauthstop:
	docker stop mariadbauth

mariadbauthrm:
	docker rm mariadbauth

createmariadb:
	docker exec -it mariadbauth mariadb --user root -psecret -e 'create database simple_bank'

mariadbmigrateup:
	migrate -path authdb/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank?parseTime=true" -verbose up

mariadbmigratedown:
	migrate -path authdb/migration -database "mysql://root:secret@tcp(localhost:3306)/simple_bank?parseTime=true" -verbose down

postgresstart:
	docker start postgres12

postgresstop:
	docker stop postgres12

postgresrm:
	docker rm postgres12

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 drop simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/galib612/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown postgresstop postgresrm postgresstart sqlc test server mock mariadb mariadbauthstart mariadbauthstop mariadbmigratedown mariadbmigrateup createmariadb mariadbauthrm