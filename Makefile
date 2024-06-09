postgresinit:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_PASSWORD=123 postgres

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres go-chat

dropdb:
	docker exec -it postgres15 dropdb go-chat

.PHONY: postgresinit postgres createdb dropdb
