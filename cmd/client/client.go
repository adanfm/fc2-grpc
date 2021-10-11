package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/adanfm/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connection to gRPC Server: %v", err)
	}

	// Quando termina de usar o connection ele fecha a conexão
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joaozinho",
		Email: "adan@adan.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joaozinho",
		Email: "adan@adan.com",
	}

	resStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := resStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())

	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "a1",
			Name:  "Adan Felipe",
			Email: "adan.grg@gmail.com",
		},
		&pb.User{
			Id:    "a2",
			Name:  "Adan Felipe 2",
			Email: "adan.grg2@gmail.com",
		},
		&pb.User{
			Id:    "a3",
			Name:  "Adan Felipe 3",
			Email: "adan.grg3@gmail.com",
		},
		&pb.User{
			Id:    "a4",
			Name:  "Adan Felipe 4",
			Email: "adan.grg4@gmail.com",
		},
		&pb.User{
			Id:    "a5",
			Name:  "Adan Felipe 5",
			Email: "adan.grg5@gmail.com",
		},
		&pb.User{
			Id:    "a6",
			Name:  "Adan Felipe 6",
			Email: "adan.grg6@gmail.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "a1",
			Name:  "Adan Felipe",
			Email: "adan.grg@gmail.com",
		},
		&pb.User{
			Id:    "a2",
			Name:  "Adan Felipe 2",
			Email: "adan.grg2@gmail.com",
		},
		&pb.User{
			Id:    "a3",
			Name:  "Adan Felipe 3",
			Email: "adan.grg3@gmail.com",
		},
		&pb.User{
			Id:    "a4",
			Name:  "Adan Felipe 4",
			Email: "adan.grg4@gmail.com",
		},
		&pb.User{
			Id:    "a5",
			Name:  "Adan Felipe 5",
			Email: "adan.grg5@gmail.com",
		},
		&pb.User{
			Id:    "a6",
			Name:  "Adan Felipe 6",
			Email: "adan.grg6@gmail.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}
