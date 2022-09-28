# Fake Dating

## Setup instructions

0. Git clone this repo - `git clone git@github.com:alexbilbie/fakedating.git`
0. `cd` into the cloned repository
0. Run `go mod download` to download the Go dependencies
0. Run `docker-compose up`
0. Connect to the MariaDB database in an SQL client and run the SQL in `setup.sql`

Database credentials:
* Host: `127.0.0.1`
* Port: `3306`
* Username: `root`
* Password: `example`

## Run the server

0. Running `go run cmd/server/main.go` will start the HTTP server on port `8000`

## Run tests

0. Run `go test ./...`