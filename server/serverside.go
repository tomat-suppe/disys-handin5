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
var HighestBid int64 = 0
var WinningBidder int32
var startTime time.Time = time.Now()
var bidder *pb.Bidder

type Server struct {
	pb.UnimplementedAuctionServer
}

func main() {
	var server = &Server{}

	startTime = time.Now()

	listener, err := TurnOnServer(server)

	//below while-loop should in theory make the server crash after 15 seconds
	//but it doesn't activate. I left it in here to show intended program,
	//how it actually works is described in the report.
	for {
		if time.Since(startTime) >= time.Second*15 {
			log.Printf("Failed to serve: %v", err)
			log.Print("Now changing server...")
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
			}
			conn.Close()
			break
		}
	}
}

// some code from previous hand-ins
func TurnOnServer(server *Server) (net.Listener, error) {
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
	}
	return listener, err
}

// regarding the requirement about bit taking amount as arg: here 'bidder'
// is a standin for this arg, as it has the field 'Bid', with the method .GetBid()
// using this, I can calculate a BidAmount for each bidding, and both server and Client
// keep track of their latest bids.
func (s *Server) Bid(ctx context.Context, in *pb.Bidder) (*pb.BidAccepted, error) {
	//below 3 lines adapted from chatgpt
	file, err := os.OpenFile("/tmp/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file")
	}
	BidAmount := in.GetBid() + 5
	if time.Since(startTime) <= time.Second*30 && BidAmount > HighestBid {
		BidAmount = in.GetBid() + 5
		message := "Bid has been accepted: " + fmt.Sprint(BidAmount)
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}
		WinningBidder = bidder.GetBidderId()

		HighestBid = BidAmount
		file.WriteString(" --- ")
		file.WriteString(fmt.Sprint(HighestBid))

		return BidAccepted, nil
	} else if time.Since(startTime) >= time.Second*30 {
		message := "Bid has been rejected, auction over."
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}

		return BidAccepted, nil
	} else if BidAmount < HighestBid || BidAmount == 0 {
		message := "Bid has been rejected as too low: " + fmt.Sprint(BidAmount)
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}

		return BidAccepted, nil
	} else {
		message := "Error while receiving bid."
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}

		return BidAccepted, nil
	}
}

// output 'ResultActionUpdate' is my name for the required output 'outcome'
func (server *Server) Result(ctx context.Context, bidder *pb.Bidder) (*pb.ResultAuctionUpdate, error) {

	if time.Since(startTime) <= time.Minute {
		message := "Auction has not ended yet, current highest bid is " + fmt.Sprint(HighestBid)
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         bidder.GetBid(),
			WinningBidderId:    bidder.GetBidderId(),
		}

		return ResultUpdate, nil
	} else {

		message := "!!! Auction has ended, highest bid was " + fmt.Sprint(HighestBid)
		ResultUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: message,
			WinningBid:         HighestBid,
			WinningBidderId:    WinningBidder,
		}

		return ResultUpdate, nil
	}
}
