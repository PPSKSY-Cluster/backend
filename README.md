# Backend for the ClusterThruster project

ClusterThruster is a webapplication to simplify the organization of sharing cluster resources.
Users will at least be able to register and subscribe to cluster resources.
This project is built in the context of the course Scalable Systems at TU berlin.

## Technologies

The server uses the fiber framework and is connected to a mongodb instance.

## Get started

To run the server use the command
```
go run main.go
```

To compile it  run
```
go build -o ./build/server main.go router.go
```