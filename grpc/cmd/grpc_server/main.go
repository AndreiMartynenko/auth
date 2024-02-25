package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/AndreiMartynenko/auth/grpc/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/*
The server struct embedding auth_v1.UnimplementedUserAPIServicesServer
is a way to ensure that your server struct implements all the methods
defined in the UserAPIServices service interface.

By embedding auth_v1.UnimplementedUserAPIServicesServer into your server struct,
your server struct implicitly implements the UserAPIServices service interface.
You can then override the methods you need to implement in your server struct
while leaving out those you don't need to implement.
This helps in keeping your code organized and ensures that you fulfill
the requirements of the gRPC service.
*/

const grpcPort = 50051

type server struct {
	auth_v1.UnimplementedUserAPIServicesServer
}

// Create
func (srv *server) Create(ctx context.Context, req *auth_v1.CreateUserRequest) (*auth_v1.CreateUserResponse, error) {

	//For testing purposes
	id := generateUniqueID()

	return &auth_v1.CreateUserResponse{
		//Id: gofakeit.Int64(),
		Id: id,
	}, nil
}

// Get
func (srv *server) Get(ctx context.Context, req *auth_v1.GetUserRequest) (*auth_v1.GetUserResponse, error) {
	log.Printf("User id %d", req.GetId())

	return &auth_v1.GetUserResponse{
		Id: req.GetId(),
		// Name:      gofakeit.Name(),
		Name:      "NEW NAME",
		Email:     gofakeit.Email(),
		Role:      auth_v1.UserRole_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

// Update
func (srv *server) Update(ctx context.Context, req *auth_v1.UpdateUserRequest) (*auth_v1.UpdateUserResponse, error) {

	updatedName := gofakeit.Name()
	updatedEmail := gofakeit.Email()
	updatedRole := auth_v1.UserRole(gofakeit.Number(0, 1))

	log.Printf("Updating user with ID %d", req.GetId())
	log.Printf("New Name: %s, New Email: %s, New Role: %v", updatedName, updatedEmail, updatedRole)

	return &auth_v1.UpdateUserResponse{}, nil

}

// Delete
func (srv *server) Delete(ctx context.Context, req *auth_v1.DeleteUserRequest) (*auth_v1.DeleteUserResponse, error) {

	err := deleteUserByID(req.GetId())
	if err != nil {
		return nil, err
	}
	return &auth_v1.DeleteUserResponse{}, nil

}

func deleteUserByID(userID int64) error {
	//Example
	return nil
}

func generateUniqueID() int64 {

	return rand.Int63n(1000)
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	/*
		Reflection in this context allows gRPC clients to query information
		about the gRPC server's services dynamically at runtime.
		It enables tools like gRPC's command-line interface (grpc_cli)
		and gRPC's web-based GUI (grpcui) to inspect the server's
		services and make RPC calls without needing to know
		the specifics of each service beforehand.
	*/
	reflection.Register(srv)
	auth_v1.RegisterUserAPIServicesServer(srv, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
