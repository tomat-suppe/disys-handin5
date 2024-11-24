package main

import (
	"log"
	"net"
	"time"

	pb "disys-handin5/proto_files"

	"google.golang.org/grpc"
)

var CurrentBid = 0
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

func (server *Server) Bid(bidder *pb.Bidder) *pb.BidAccepted {
	if time.Since(startTime) < 500 {
		CurrentBid = CurrentBid + 5
		message := "Bid has been accepted: " + string(CurrentBid)
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}
		return BidAccepted
	} else {
		message := "Bid has been rejected. Auction is over."
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}
		return BidAccepted
	}
}

func (server *Server) Result(bidder *pb.Bidder) *pb.ResultAuctionUpdate {
	if time.Since(startTime) < 500 {
		message := "Auction has not ended yet, current highest bid is " + string(CurrentBid)
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         int64(CurrentBid),
			WinningBidderId:    0,
		}
		return ResultUpdate
	} else {
		message := "!!! Auction has ended, highest bid was " + string(CurrentBid)
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         int64(CurrentBid),
			WinningBidderId:    0,
		}
		return ResultUpdate
	}
}
