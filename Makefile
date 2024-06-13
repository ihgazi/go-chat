postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_PASSWORD=123 postgres

postgres:
	docker exec -it postgres15 psql -U postgres

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres go-chat

dropdb:
	docker exec -it postgres15 dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:123@localhost:5432/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb
