DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres: 
	docker run --name postgres-golang --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb: 
	docker exec -it postgres-golang createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it postgres-golang dropdb simple_bank

migratecreate: 
	@migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown: 
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1: 
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs: 
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc: 
	sqlc generate

test: 
	go test -v -cover ./...

server: 
	go run main.go

mock: 
	mockgen -destination db/mock/store.go -package mockdb github.com/kien-ltt/simplebank/db/sqlc Store

proto: 
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

evans: 
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createdb dropdb migratecreate migrateup migratedown migrateup1 migratedown1 db_docs db_schema sqlc test server mock proto evans
	