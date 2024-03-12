postgres:
	 docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -d postgres
createdb:
	docker exec -it postgres createdb --username=root --owner=root simplebank
dropdb:
	docker exec -it postgres dropdb simplebank
migrate-up:
	migrate -path db/migration -database "postgresql://root:mysecretpassword@localhost:5432/simplebank?sslmode=disable" -verbose up
migrate-down:
	migrate -path db/migration -database "postgresql://root:mysecretpassword@localhost:5432/simplebank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v --cover ./...
.PHONY: postgres  createdb dropdb migrateup migratedown