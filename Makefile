.PHONY: run-server run-client db-start migrate-up migrate-down setup

run-server:
	go run *.go

run-client:
	cd client && npm run start

db-start:
	docker-compose -f db/docker-compose.yml up

migrate-up:
	migrate -database 'mysql://root:@tcp(127.0.0.1:3307)/sudoku' -path ./db/migrations up

migrate-down:
	migrate -database 'mysql://root:@tcp(127.0.0.1:3307)/sudoku' -path ./db/migrations down

setup:
	go get github.com/go-sql-driver/mysql
	go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq
	go build -tags 'mysql' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli
	-go get .
	cd client && npm install