package main

import (
	"bufio"
	"log"
	"net"
	"time"

	pb "disys-handin5/protofiles"

	"google.golang.org/grpc"
)

var time time

type Server struct {
	pb.UnimplementedAuctionServer
	port int
}

var server = &Server{
	port: 50000,
}

var backupserver = &Server{
	port: 50000,
}


func main() {
	go TurnOnServer(server)
	go server.Bid(args)

	time := time.Now()
}

func TurnOnServer(server *Server) {
	//some code from previous hand-ins
	listener, err := net.Listen("tcp", server.port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Printf("Server now listening on: %s", server.port)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterAuctionServer(grpcServer)

	log.Printf("Server is running on : %s ...", server.port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		TurnOnServer(backupserver)
	}
	if time.Since(time) > 200 {
		log.Fatalf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		TurnOnServer(backupserver)
	}
}

func (server *Server) Bid (ctx context.Context, bidder *pb.Bidder) (stream *pb.AuctionUpdate, error){
	//for{
	input, err := stream.Recv()
	if err != nil {
		log.Fatalf("Server cannot update auction")
	}
	//}
}

func (server *Server) Result (ctx context.Context, bidder *pb.Bidder) (stream *pb.AuctionUpdate, error){
	//for{
	input, err := stream.Recv()
	if err != nil {
		log.Fatalf("Server cannot update auction")
	}
	//}
}
