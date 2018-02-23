package main

import (
	"encoding/json"
	"grpc_api/proto"
	"log"
	"net/http"

	context "golang.org/x/net/context"

	"google.golang.org/grpc"
)

type Config struct {
	Grpcl *grpc.ClientConn
}

func createUser(w http.ResponseWriter, r *http.Request, config Config) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Not a form")
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	name := r.PostFormValue("name")

	client := proto.NewUserServiceClient(config.Grpcl)
	user, err := client.AddUser(context.Background(), &proto.AddUserRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})
	if err != nil {
		log.Printf("failed with %s", err)
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to create json %s", err)
	}

	w.Write(b)
}

func loginUser(w http.ResponseWriter, r *http.Request, config Config) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Not a form")
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	client := proto.NewUserServiceClient(config.Grpcl)
	user, err := client.LoginUser(context.Background(), &proto.LoginUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Printf("failed with %s", err)
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to create json %s", err)
	}

	w.Write(b)
}

func getUser(w http.ResponseWriter, r *http.Request, config Config) {
	authHeader := r.Header.Get("Authorization")
	userID := r.URL.Query().Get("id")

	client := proto.NewUserServiceClient(config.Grpcl)
	user, err := client.GetUser(context.Background(), &proto.GetUserRequest{
		Headers: map[string]string{"authorization": authHeader},
		Params:  map[string]string{"id": userID},
	})
	if err != nil {
		log.Printf("failed with %s", err)
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to create json %s", err)
	}

	w.Write(b)
}

func main() {
	grpClient, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to grpc server at port 7777")
	}
	config := Config{Grpcl: grpClient}

	http.HandleFunc("/create_user", func(w http.ResponseWriter, r *http.Request) {
		createUser(w, r, config)
	})
	http.HandleFunc("/login_user", func(w http.ResponseWriter, r *http.Request) {
		loginUser(w, r, config)
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		getUser(w, r, config)
	})

	http.ListenAndServe(":8000", nil)
}
