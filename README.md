### grpc_api

Unary for now

#### Requirements
- golang
- protoc (optional)
- google.golang.org/grpc

#### How to run
- go run server/main.go
- go run server/client.go

#### Curl requests
- `curl localhost:8000/create_user -H 'Content-Type: application/x-www-form-urlencoded' -d 'email=k@b.c&name=kbc&password=abcd1234`
- `curl localhost:8000/login_user -H 'Content-Type: application/x-www-form-urlencoded' -d 'email=k@b.c&password=abcd1234'`
- `curl localhost:8000/user?id=1 -H 'Authorization: <Token>`

