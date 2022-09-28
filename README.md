# Fake Dating

## Setup instructions

1. Git clone this repo - `git clone git@github.com:alexbilbie/fakedating.git`
2. `cd` into the cloned repository
3. Run `go mod download` to download the Go dependencies
4. Run `docker-compose up`
5. Connect to the MariaDB database in an SQL client and run the SQL in `setup.sql`

Database credentials:
* Host: `127.0.0.1`
* Port: `3306`
* Username: `root`
* Password: `example`

## Run the server

1. Running `go run cmd/server/main.go` will start the HTTP server on port `8000`

## Run tests

1. Run `go test ./...` from the root - requires Docker to be running

* `pkg/middleware/auth_test.go` is an example of testing an HTTP endpoint
* `pkg/repository/user_test.go` is an example of testing a database using a Docker container
* `pkg/util/*_test.go` shows simple unit tests

## API

### Create a user

Request: 

```
POST /user/create
```

Response:

```
HTTP/1.1 201 Created
Content-Type: application/json

{
    "ID": "2FOKFW4brtIY8Mf7sQJGpwM4GLa",
    "Email": "mlhjypj@yiwciki.ru",
    "Name": "Cordelia VonRueden",
    "Gender": 2,
    "Age": 25,
    "Location": {
        "Latitude": 51.56771049736385,
        "Longitude": 0.041990415849927736
    },
    "Password": "Pld!!QF12SaN8Vl"
}
```

### Login

Request:

```
POST /login
Content-Type: application/json

{
    "Email": "mlhjypj@yiwciki.ru",
    "Password": "Pld!!QF12SaN8Vl"
}
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "User": {
        "ID": "2FOKFW4brtIY8Mf7sQJGpwM4GLa",
        "Email": "mlhjypj@yiwciki.ru",
        "Name": "Cordelia VonRueden",
        "Gender": 2,
        "Age": 25,
        "Location": {
            "Latitude": 51.56771049736385,
            "Longitude": 0.041990415849927736
        }
    },
    "Token": "2FOKFt5HTqs6sg4ncxnRV65KPZ5"
}
```

### Get profiles

Request:

```
GET /profiles
Authorization: 2FOKFt5HTqs6sg4ncxnRV65KPZ5
```

Available query parameters:

* `latitude` (float64) - requires `longitude` to be set
* `longitude` (float64) - requires `latitude` to be set
* `radius` (uint) - KM radius - default 1
* `age_lower` (uint) - lower age bound - min 18
* `age_upper` (uint) - upper age bound - max 99
* `offet` (uint) - pagination offset

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "Matches": [
        {
            "ID": "2FOKcCVfQA8kAclB4o4VqTiatDQ",
            "Email": "xylttoa@wsjncam.info",
            "Name": "Bernhard Rohan",
            "Gender": 2,
            "Age": 44,
            "Location": {
                "Latitude": 51.47817688295783,
                "Longitude": 0.027873615525135505
            }
        }
    ]
}
```

### Swipe

Request:

```
POST /swipe HTTP/1.1
Authorization: 2FOKFt5HTqs6sg4ncxnRV65KPZ5

{
    "Recipient":"2FOKcCVfQA8kAclB4o4VqTiatDQ", 
    "Matched":true // or false if user selected no
}
```

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "MutualMatch": false
}
```
