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
var HighestBid int64 //arbitrary artificial number for the sake of running program.
// see report for further explanation
var startTime time.Time = time.Now()
var bidder *pb.Bidder
var serverCrashed bool = false

type Server struct {
	pb.UnimplementedAuctionServer
}

func main() {
	var server = &Server{}

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

// some code from previous hand-ins
func TurnOnServer(server *Server) {
	//here I would implement logic to get the latest highest bid.
	//as it is right now, every bidder keeps track of their own highest bid,
	//but there could be an issue with a bid too small being accepted, even
	//if the logs look 'normal, because this server does not know what the previous highest bid was.

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

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)

	}

	file, err := os.OpenFile("/tmp/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file")
	}
	stat, err := file.Stat()
	Bid, err := file.ReadAt(make([]byte, 1), stat.Size()-1)
	BidAmount = int64(Bid)
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
	if time.Since(startTime) <= time.Minute && BidAmount > HighestBid {
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
	} else if time.Since(startTime) >= time.Minute {
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
