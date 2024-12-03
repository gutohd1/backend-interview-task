package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "app/explore_service_protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:9000 : %v", err)
	}

	defer conn.Close()

	c := pb.NewExploreServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var command, actorId, recipientId, page string
	var like bool
	flag.StringVar(&command, "function", "", "-function=ListLikedYou, ListNewLikedYou, CountLikedYou and PutDecision")
	flag.StringVar(&actorId, "actor", "1", "-actor=1 | id to call specific actor user")
	flag.StringVar(&recipientId, "recipient", "1", "-recipient=1 | id to call specific recipient user")
	flag.BoolVar(&like, "like", true, "-like=false | Can only be used on PutDecision")
	flag.StringVar(&page, "page", "1", "-page=1 | page number to paginate likes")
	flag.Parse()

	switch command {
	case "ListLikedYou":
		request := &pb.ListLikedYouRequest{
			RecipientUserId: recipientId,
			PaginationToken: &page,
		}

		listLikedYou, err := c.ListLikedYou(ctx, request)
		if err != nil {
			log.Fatalf("error calling function ListLikedYou: %v", err)
		}
		log.Printf("response from server %+v", listLikedYou)
	case "ListNewLikedYou":
		fmt.Println("Calling ListNewLikedYou function")
		request := &pb.ListLikedYouRequest{
			RecipientUserId: recipientId,
			PaginationToken: &page,
		}

		ListNewLikedYou, err := c.ListNewLikedYou(ctx, request)
		if err != nil {
			log.Fatalf("error calling function ListNewLikedYou: %v", err)
		}
		log.Printf("response from server %+v", ListNewLikedYou)
	case "CountLikedYou":
		fmt.Println("Calling CountLikedYou function")
		request := &pb.CountLikedYouRequest{
			RecipientUserId: recipientId,
		}

		CountLikedYou, err := c.CountLikedYou(ctx, request)
		if err != nil {
			log.Fatalf("error calling function CountLikedYou: %v", err)
		}
		log.Printf("response from server %+v", CountLikedYou)
	case "PutDecision":
		fmt.Println("Calling PutDecision function")
		request := &pb.PutDecisionRequest{
			ActorUserId:     actorId,
			RecipientUserId: recipientId,
			LikedRecipient:  like,
		}

		PutDecision, err := c.PutDecision(ctx, request)
		if err != nil {
			log.Fatalf("error calling function PutDecision: %v", err)
		}
		log.Printf("response from server %+v", PutDecision)
	default:
		fmt.Println("No command were selected")
	}

}
