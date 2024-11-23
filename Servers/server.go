package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "disys-handin5/protofiles"

	"google.golang.org/grpc"
)

var startTime time.Time

type Server struct {
	pb.UnimplementedAuctionServer
	port string
}

var server = &Server{
	port: "localhost:50000",
}

var backupserver = &Server{
	port: "localhost:50000",
}

func main() {
	TurnOnServer(server)
	//go server.Bid()

	startTime = time.Now()
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

	pb.RegisterAuctionServer(grpcServer, server.UnimplementedAuctionServer)

	log.Printf("Server is running on : %s ...", server.port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		TurnOnServer(backupserver)
	}
	if time.Since(startTime) > 200 {
		log.Fatalf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		TurnOnServer(backupserver)
	}
}

func (server *Server) Bid(ctx context.Context, bidder *pb.Bidder) (stream pb.Auction_BidClient) {
	//for{
	input, err := stream.Recv()
	if err != nil {
		log.Fatalf("Server stopped working")
	} else {
		log.Printf("Bidder #%v has bid %v", input.HighestBidderId, input.HighestBid)
	}
	return nil
	//}
}

func (server *Server) Result(ctx context.Context, bidder *pb.Bidder) (stream pb.Auction_ResultClient) {
	//for{
	input, err := stream.Recv()
	if err != nil {
		log.Fatalf("Server stopped working")
	} else {
		log.Printf("Current highest bid is %v by bidder #%v", input.WinningBid, input.WinningBidderId)
	}
	return nil
	//}
}
