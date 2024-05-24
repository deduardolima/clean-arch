createmigration:
	migrate create -ext=sql -dir=/root/migrations -seq init

migrate:
	migrate -path=/root/migrations -database "mysql://root:root@tcp(mysql:3306)/orders" -verbose up

migratedown:
	migrate -path=/root/migrations -database "mysql://root:root@tcp(mysql:3306)/orders" -verbose down

.PHONY: migrate migratedown createmigration

