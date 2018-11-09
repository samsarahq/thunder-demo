# Thundoku

This directory contains a basic example of Thunder. The code is organized as follows:

- All server code can be found in the `main.go` file.
- The `db/` directory contains MySQL schema and configuration files.
- The `client/` directory contains a JavaScript client.

## Dependencies

This example requires [Docker](https://docker.com), [Go](https://golang.org/),
and [Node](https://nodejs.org/).

Make sure the project is cloned into your `$GOPATH` and that `$GOPATH` is properly set.

## Quickstart

To quickly get started, the individual steps described later in the readme have been combined into a makefile.

* To get set up, run `make setup`. This will install all server and client dependencies.
* To start the database in the foreground, run `make db-start`. Run `make migrate-up` to setup the database's schema.
* To start the server, run `make run-server`.
* To start the client, run `make run-client`.

## Running the database

To start the database, run `docker-compose -f db/docker-compose.yml up` to
start a MySQL server on port 3307, properly configured for use with Thunder. 

Then, install [migrate](https://github.com/golang-migrate/migrate/tree/master/cli) using:
```
$ go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq
$ go build -tags 'mysql' -o /usr/local/bin/migrate github.com/golang-migrate/migrate/cli
```
If you encounter any errors such as `cannot find package "github.com/go-sql-driver/mysql"`,
you will need to `go get` each of the packages e.g. `go get github.com/go-sql-driver/mysql`.

Then run the following to set-up the database's schema:
```
migrate -database 'mysql://root:@tcp(127.0.0.1:3307)/sudoku' -path ./db/migrations up
```
Now you can access the database with `mysql -h 127.0.0.1 --port=3307 -uroot
sudoku`. Try inserting a new message by running
`INSERT INTO games (state, data) VALUES (0, "testdata");`

## Running the server

To run the server, first install the server's dependencies using
`go get .`.
Then, run `go run *.go` to start the server.

## Running the client

To run the client, first install the client's dependencies using `npm install`.
Then, run `npm run start` to start the client. You can access the basic client
by going to `http://localhost:3000`. To run your own queries, access the
GraphiQL client at `http://localhost:3000/graphiql`. One example GraphQL query
to fetch all messages is as follows:
```
{
  messages {
    id
    text
  }
}
```
