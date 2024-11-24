package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "disys-handin5/proto_files"

	"google.golang.org/grpc"
)

var Bid int64 = 0
var WinningBidder int32
var startTime time.Time
var bidder *pb.Bidder

type Server struct {
	pb.UnimplementedAuctionServer
	port string
}

func main() {
	var server = &Server{
		port: "localhost:50000",
	}
	TurnOnServer(server)

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

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	pb.RegisterAuctionServer(grpcServer, &Server{})

	log.Printf("Server is running on : %s ...", server.port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		backupserver := &Server{
			port: "localhost:50000",
		}
		TurnOnServer(backupserver)
	}
	if time.Since(startTime) > 200 {
		log.Fatalf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		backupserver := &Server{
			port: "localhost:50000",
		}
		TurnOnServer(backupserver)
	}
	//go server.Bid(ctx, bidder)

}

func (s *Server) Bid(ctx context.Context, in *pb.Bidder) (*pb.BidAccepted, error) {
	if time.Since(startTime) < 500 {
		Bid = bidder.GetBid() + 5
		bidder.Bid = Bid
		message := "Bid has been accepted: " + fmt.Sprint(Bid)
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}
		WinningBidder = bidder.GetBidderId()
		return BidAccepted, nil
	} else {
		message := "Bid has been rejected. Auction is over."
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}
		return BidAccepted, nil
	}
}

func (server *Server) Result(bidder *pb.Bidder) (*pb.ResultAuctionUpdate, error) {
	if time.Since(startTime) < 500 {
		message := "Auction has not ended yet, current highest bid is " + fmt.Sprint(bidder.GetBid())
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         bidder.GetBid(),
			WinningBidderId:    bidder.GetBidderId(),
		}
		return ResultUpdate, nil
	} else {

		message := "!!! Auction has ended, highest bid was " + string(Bid)
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         Bid,
			WinningBidderId:    WinningBidder,
		}
		return ResultUpdate, nil
	}
}

func (server *Server) SendBid(ctx context.Context, data *pb.Bidder) (*pb.Bidder, error) {
	bid := data.GetBid()
	bidderid := data.GetBidderId()
	bidder = &pb.Bidder{
		BidderId: bidderid,
		Bid:      bid,
	}

	return bidder, nil
}
