# Backend for the ClusterThruster project

ClusterThruster is a web application to simplify the organization of sharing cluster resources.
Users will at least be able to register and subscribe to cluster resources.
This project is built in the context of the course Scalable Systems at TU Berlin.

## Technologies

The [go](https://go.dev/) server uses the [fiber framework](https://docs.gofiber.io/) and is connected to a mongodb instance, using the [mongo driver](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo?utm_source=godoc).
Authentication/authorization is done using jwt and bcrypt.

## Get started

To just run the backend in docker use
```
docker-compose up
```
Make sure that you adapt your .env as described below

To run the server use the command
```
go run main.go
```

To compile it  run
```
go build -o ./build/server main.go
```
-o to specify output directory

To run tests execute
```
docker-compose -f ./tests/docker-compose.yml up -d
```
-For running a local instance of mongodb and the fake-smtp-server
```
go test -v ./tests -tags test
```
-v for verbosity
./tests because that's the directory where they live
-tags test because the ability to drop the entire database is only compiled for testing

For debugging we recommend using [delve](https://github.com/go-delve/delve)
```
dlv debug github.com/PPSKSY-Cluster/backend
```

## Docs

Docs are provided via [swagger](https://github.com/swaggo/swag) and located under the route `/api/docs` or under our [github pages](https://ppsksy-cluster.github.io/backend/).
You may regenerate them by running 
```
swag init
```
but there is also a GitHub action running to update them and open pull requests.

See [declarative comments](https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html) if you want to add some information.


## Configuration and secrets

Currently, all that is located in your .env in the root directory.
Possible fields are :

| key                   | example value                  | description                                                                                                                                         |
|-----------------------|--------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| DB_NAME               | test                           | The name of your database                                                                                                                           |
| MONGODB_URI           | mongodb://localhost:27017/test | The connection string to your database<br/>(change 'localhost' to 'mongo_db' when using docker)                                                     |
| PORT                  | 8080                           | The port to start the server on                                                                                                                     |
| CLIENT_URL            | http://localhost:3000          | The URL the client is requesting from                                                                                                               |
| BCRYPT_COST           | 10                             | The Cost for bcrypt pw generation (>= 10)                                                                                                           |
| RSA_KEY_SIZE          | 2048                           | The key size used when generating rsa key pairs                                                                                                     |
| RELAY_HOST            | mail.gmx.net                   | The name of the relay host used for transferring messages sent by the smtp client <br/> Change to 'localhost' for testing                           |
| RELAY_PORT            | 587                            | The port of the relay host  <br/> Change to 5025 for testing                                                                                        |
| RELAY_USERNAME        | cluster-thruster@gmx.net       | The username of the account used on the relay host <br/> Omit for testing                                                                           |
| RELAY_PASSWORD        | changeme                       | The password of the account used on the relay host <br/> Omit for testing                                                                           |
| NOTIFICATION_INTERVAL | 3600                           | The interval in which the availability of wanted resources is checked (in seconds)<br/>for testing use smaller intervals                            |

For tests a seperate .env is needed in the corresponding folder `./tests`.

## Short note on mail
We've used the mail functionality with [gmx](https://www.gmx.net/) as a relay host, but all other relays supporting username and password based authentication should be fine.
<br/> For testing purposes we use [fake-smtp-server](https://github.com/gessnerfl/fake-smtp-server).

## Repository structure
```
> api        // router and route handlers package
> auth       // authentication/authorization package
> build      // compiled files after running build command
> db         // the db package with the various interfaces
> dist       // needed by index.html
> docs       // docs as generated by swagger
> mail       // helper functions and templates for mail functionality
> tests      // api tests
.env
.gitignore
docker-compose.yml
Dockerfile
go.mod
go.sum
index.html   // for github pages
LICENSE
main.go      // server entry point, initalizes as starts everything
```
