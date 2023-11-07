migration.create:
	# create new migration
	migrate create -ext sql -dir database/migrations -seq $(name)

migration.up:
	# execute migration up
	migrate --path database/migrations -database "postgres://root:secret@localhost:5432/debozero_db?sslmode=disable" up
	migrate --path database/migrations -database "postgres://root:secret@localhost:5433/debozero_db?sslmode=disable" up

migration.down:
	# execute rollback migration
	migrate --path database/migrations -database "postgres://root:secret@localhost:5432/debozero_db?sslmode=disable" down
	migrate --path database/migrations -database "postgres://root:secret@localhost:5433/debozero_db?sslmode=disable" down

postgres.up:
	# create postgres db server
	docker-compose up -d

postgres.down:
	# delete postgres db server
	docker-compose down

postgres.db.up:
	# create db
	docker exec -it debozero_postgres createdb --username=root --owner=root debozero_db
	docker exec -it debozero_postgres_live createdb --username=root --owner=root debozero_db

postgres.db.down:
	# drop db
	docker exec -it debozero_postgres dropdb --username=root debozero_db
	docker exec -it debozero_postgres_live dropdb --username=root debozero_db

sqlc.generate:
	# generate models using sqlc
	docker run --rm -v "D:\work\github.com\debozero-backend:/src" -w /src sqlc/sqlc generate

test:
	go test -v -cover ./...