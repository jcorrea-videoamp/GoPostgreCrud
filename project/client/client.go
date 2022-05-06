package main

import (
	"context"
	"log"

	"github.com/jcorrea-videoamp/GoPostgreCrud/project/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewOrderServiceClient(conn)
	orders, err := client.ListOrders(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalln("client failed to fetch orders")
		return
	}
	if len(orders.Orders) > 0 {
		log.Println("Data:", orders.Orders[0].Customer)
	}
	log.Println("Here")
}
