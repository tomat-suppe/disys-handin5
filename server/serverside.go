package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/tomat-suppe/disys-handin5/proto_files"

	"google.golang.org/grpc"
)

var BidAmount int64 = 0
var WinningBidder int32
var startTime time.Time = time.Now()
var bidder *pb.Bidder

type Server struct {
	pb.UnimplementedAuctionServer
}

func main() {
	var server = &Server{}

	//startTime = time.Now()

	TurnOnServer(server)

}

func TurnOnServer(server *Server) {
	//some code from previous hand-ins
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Printf("Server now listening on: localhost:50000")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuctionServer(grpcServer, server)

	log.Printf("Server is running on : localhost:50000")
	if err := grpcServer.Serve(listener); err != nil {
		log.Printf("Failed to serve: %v", err)
		/*log.Print("Now changing server...")
		backupserver := &Server{}
		TurnOnServer(backupserver)*/
	}
	if time.Since(startTime) >= time.Minute {
		log.Printf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
		}
		conn.Close()
		/*CloseConnection(conn)
		log.Fatalf("Failed to serve: %v", err)
		log.Print("Now changing server...")
		backupserver := &Server{}
		TurnOnServer(backupserver)*/
	}

	//go server.Bid(ctx, bidder)
}

func (s *Server) Bid(ctx context.Context, in *pb.Bidder) (*pb.BidAccepted, error) {
	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND, 0666) //doesn't work yet
	if err != nil {
		log.Printf("Failed to open file")
	}
	if time.Since(startTime) <= time.Minute {
		BidAmount = BidAmount + 5
		//bidder.Bid = BidAmount
		message := "Bid has been accepted: " + fmt.Sprint(BidAmount)
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}
		WinningBidder = bidder.GetBidderId()
		file.WriteString(message)
		return BidAccepted, nil
	} else {
		message := "Bid has been rejected. Auction is over."
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}
		file.WriteString(message)
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
