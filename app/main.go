package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"app/database"
	pb "app/explore_service_protos"
	"app/handlers"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()
	lis, err := net.Listen("tcp", ":"+os.Getenv("LISTEN_PORT"))
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := sql.Open("mysql", connection)

	if err != nil {
		panic(err.Error())
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExploreServiceServer(grpcServer, &handlers.Server{
		DatabaseReader: database.NewDatabaseReader(db),
		DatabaseWriter: database.NewDatabaseWriter(db),
	})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Feiled to serve %s", err)
	}
}
