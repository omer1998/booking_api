up_migration_test_db:
	migrate -path migrations -database postgres://postgres:12345678@127.0.0.1:5444/booking?sslmode=disable up

up_migration_dev_db:
	migrate -path migrations -database postgres://postgres:12345678@localhost:5433/booking?sslmode=disable up

down_migration_test_db:
	migrate -path migrations -database postgres://postgres:12345678@127.0.0.1:5444/booking?sslmode=disable down
down_migration_dev_db:
	migrate -path migrations -database postgres://postgres:12345678@localhost:5433/booking?sslmode=disable down

create_migration:
	migrate create -ext sql -dir migrations -seq $(name)