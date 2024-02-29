package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AndreiMartynenko/auth/grpc/pkg/auth_v1"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

const (
	dbDSN     = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
	grpcPort  = 50051
	dbTimeout = time.Second * 3
)

type server struct {
	db *sql.DB
	auth_v1.UnimplementedUserAPIServicesServer
}

func main() {
	log.Println("Starting authentication service")

	ctx := context.Background()

	// Create a connection to the database
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()

	reflection.Register(srv)
	auth_v1.RegisterUserAPIServicesServer(srv, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

// Create
func (srv *server) Create(ctx context.Context, req *auth_v1.CreateUserRequest) (*auth_v1.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Insert the new user into the database
	res := srv.db.QueryRowContext(ctx, "INSERT INTO users (name, email, password, password_confirmed, role) VALUES ($1, $2, $3, $4, $5) RETURNING id", req.Name, req.Email, req.Password, req.PasswordConfirmed, req.Role)
	// Get the ID of the newly created user
	var id int64
	err := res.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to get id of newly created user: %v", err)
	}

	// if res.RowsAffected() > 0 {
	// 	id = res.RowsAffected()
	// } else {

	return &auth_v1.CreateUserResponse{
		Id: id,
	}, nil

}

// Get
func (srv *server) Get(ctx context.Context, req *auth_v1.GetUserRequest) (*auth_v1.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Query the user from the database
	row := srv.db.QueryRowContext(ctx, "SELECT name, email, role, created_at, updated_at FROM users WHERE id = $1", req.GetId())

	// Scan the result into variables
	var name string
	var email string
	var role int32
	var createdAt time.Time
	var updatedAt time.Time
	err := row.Scan(&name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
		}
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	return &auth_v1.GetUserResponse{
		Id:        req.GetId(),
		Name:      name,
		Email:     email,
		Role:      auth_v1.UserRole(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

// Update
func (srv *server) Update(ctx context.Context, req *auth_v1.UpdateUserRequest) (*auth_v1.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Update the user in the database
	result, err := srv.db.ExecContext(ctx, "UPDATE users SET name = $1, email = $2, role = $3 WHERE id = $4", req.GetName().GetValue(), req.GetEmail().GetValue(), req.GetRole(), req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
	}

	log.Printf("Updating user with ID %d", req.GetId())
	log.Printf("New Name: %s, New Email: %s, New Role: %v", req.GetName().GetValue(), req.GetEmail().GetValue(), req.GetRole())

	return &auth_v1.UpdateUserResponse{}, nil
}

// Delete
func (srv *server) Delete(ctx context.Context, req *auth_v1.DeleteUserRequest) (*auth_v1.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Delete the user from the database
	result, err := srv.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
	}

	return &auth_v1.DeleteUserResponse{}, nil
}
