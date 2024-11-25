package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/tomat-suppe/disys-handin5/proto_files"

	"google.golang.org/grpc"
)

var BidAmount int64 = 0
var WinningBidder int32
var startTime time.Time = time.Now()
var bidder *pb.Bidder
var serverCrashed bool = false

type Server struct {
	pb.UnimplementedAuctionServer
}

func main() {
	var server = &Server{}
	//startTime = time.Now()

	for {
		ListenForServerCrash(server)
		if serverCrashed == true {

			break
		}
	}
	log.Printf("Backupserver is connecting...")
	TurnOnServer(server)

}

func ListenForServerCrash(server *Server) {
	//below 3 lines adapted from https://stackoverflow.com/questions/56336168/golang-check-tcp-port-open
	timeout := time.Second
	_, err := net.DialTimeout("tcp", "localhost:50000", timeout)
	if err != nil {
		serverCrashed = true
	}
}

func TurnOnServer(server *Server) {
	//some code from previous hand-ins
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		log.Fatalf("Backupserver failed to listen: %v", err)
	} else {
		log.Printf("Backupserver now listening on: localhost:50000")
	}
	defer listener.Close()
	grpcServer := grpc.NewServer()

	pb.RegisterAuctionServer(grpcServer, server)

	log.Printf("Backupserver is running on : localhost:50000")
	/*for {

	}*/
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)

	}
	//go server.Bid(ctx, bidder)
}

func (s *Server) Bid(ctx context.Context, in *pb.Bidder) (*pb.BidAccepted, error) {
	if time.Since(startTime) <= time.Minute {
		BidAmount = BidAmount + 5
		//bidder.Bid = BidAmount
		message := "Bid has been accepted: " + fmt.Sprint(BidAmount)
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

func (server *Server) Result(ctx context.Context, bidder *pb.Bidder) (*pb.ResultAuctionUpdate, error) {
	if time.Since(startTime) <= time.Minute*2 {
		message := "Auction has not ended yet, current highest bid is " + fmt.Sprint(BidAmount)
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         bidder.GetBid(),
			WinningBidderId:    bidder.GetBidderId(),
		}
		return ResultUpdate, nil
	} else {

		message := "!!! Auction has ended, highest bid was " + string(BidAmount)
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         BidAmount,
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
