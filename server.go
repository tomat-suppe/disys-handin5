package main

import (
	"bufio"
	"log"
	"net"

	pb "disys-handin5/protofiles"

	"google.golang.org/grpc"
)

var Scanner bufio.Scanner
var HighestBid int = 0
var FinalBid int = 0
var TimeSteps int = 0 //maybe this needs to be an actual time component

type Server struct {
	pb.UnimplementedAuctionServer
	port int
}

func main() {
	for i := 0; i < 3; i++ {
		portNo := 50000 + i
		server := &Server{
			port: portNo,
		}
		go TurnOnServer(server)
		//go server.Bid(args)
	}

	//listen for bids
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
	}
}

/*func (server *Server) Bid (ctx context.Context, bidder *pb.Bidder) (*pb.AuctionUpdate, error){
}*/
