This directory contains the backend for our thunder-demo app. The code is organized as follows:

- All server code can be found in the `server/main.go` file.
- The `db/` directory contains MySQL schema and configuration files.

## Dependencies

This example requires [Docker](https://docker.com), [Go](https://golang.org/),
and [Node](https://nodejs.org/).

## Running the database

To start the database, run `docker-compose -f db/docker-compose.yml up` to
start a MySQL server on port 3307, properly configured for use with Thunder. 

Then, install [migrate](https://github.com/golang-migrate/migrate/tree/master/cli) using:
```
$ go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq
$ go build -tags 'mysql' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli
```
Then run
```
migrate -database mysql://root:@tcp(127.0.0.1:3307)/github -path ./db/migrations up
```
to set-up the database's schema.

Now you can access the database with `mysql -h 127.0.0.1 --port=3307 -uroot
github`. Try inserting a new message by running
`INSERT INTO repos (id, full_name, api_json) VALUES (1, "samsarahq/thunder", "");`

## Running the server

To run the server, first install the server's dependencies using
`go get .`.
Then, run `go run main.go` to start the server.
