package main

import (
	"context"
	"log"
	"time"

	"github.com/AndreiMartynenko/auth/grpc/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	userID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed toconnect to server: %v", err)
	}
	defer conn.Close()

	c := auth_v1.NewUserAPIServicesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &auth_v1.GetUserRequest{Id: userID})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	user := r
	log.Printf("User info:\nID: %d\nName: %s\nEmail: %s\nRole: %s\nCreated At: %v\nUpdated At: %v\n",
		user.GetId(), user.GetName(), user.GetEmail(), user.GetRole().String(),
		user.GetCreatedAt().AsTime(), user.GetUpdatedAt().AsTime())

	//log.Printf(color.RedString("User info:\n"), color.GreenString("%+v", r.GetId()))

}
