# Backend for the ClusterThruster project

ClusterThruster is a webapplication to simplify the organization of sharing cluster resources.
Users will at least be able to register and subscribe to cluster resources.
This project is built in the context of the course Scalable Systems at TU berlin.

## Technologies

The server uses the fiber framework and is connected to a mongodb instance.
Authentication/authorization is done using jwt and bcrypt.

## Get started

To run the server use the command
```
go run main.go router.go
```

To compile it  run
```
go build -o ./build/server main.go router.go
```
-o to specify output directory

To run tests execute
```
go test -v ./...
```
-v for verbosity
./... to test all subdirectories

## Docs

Docs are provided via swagger and located under the route /api/docs.
You may regenerate them by running 
```
swag init
```

## Configuration and secrets

Currently all that is located in your .env in the root directory.
Possible fields are :

| key        | example value                    | description                              |
|------------|----------------------------------|------------------------------------------|
|DB_NAME     | test                             | The name of your database                |
|MONGODB_URI | mongodb://localhost:27017/test   | The connection string to your database   |
|PORT        | 8080                             | The port to start the server on          |
|CLIENT_URL  | http://localhost:3000            | The URL the client is requesting from    |
|BCRYPT_COST | 10                               | The Cost for bcrypt pw generation (>= 10)|