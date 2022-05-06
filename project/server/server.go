package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jcorrea-videoamp/GoPostgreCrud/project/proto"
	"github.com/jcorrea-videoamp/GoPostgreCrud/project/repository"
	"github.com/jcorrea-videoamp/GoPostgreCrud/project/service"
	"google.golang.org/grpc"
)

func main() {
	driver := "postgres"
	url := "postgres://frtzcnqy:pYvsWxUKNQhG6xtFFqAj6sdTZdoc0lvB@chunee.db.elephantsql.com/frtzcnqy"
	db, err := repository.ConnectDb(driver, url)
	if err != nil {
		log.Fatalln("failed to connect to database")
		return
	}
	repo, err := repository.NewRepository(db)
	if err != nil {
		log.Fatalln("failed to start repository")
		return
	}
	srv, err := service.NewOrderService(repo)
	if err != nil {
		log.Fatalln("failed to start service")
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterOrderServiceServer(server, srv)
	server.Serve(lis)
}
