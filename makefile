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
migrate-up1:
	migrate -path db/migration -database "postgresql://root:mysecretpassword@localhost:5432/simplebank?sslmode=disable" -verbose up 1
migrate-down1:
	migrate -path db/migration -database "postgresql://root:mysecretpassword@localhost:5432/simplebank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v --cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb  --destination db/mock/store.go nebil/golang/db/sqlc Store
.PHONY: postgres  createdb dropdb migrateup migratedown server mock migrateup1 migratedown1