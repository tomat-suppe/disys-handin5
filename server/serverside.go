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

var BidAmount int64
var HighestBid int64
var WinningBidder int32
var startTime time.Time
var bidder *pb.Bidder

type Server struct {
	pb.UnimplementedAuctionServer
}

func main() {
	var server = &Server{}

	startTime = time.Now()

	TurnOnServer(server)

	//keeps server running despite the go-routine .Serve(...) in TurnOnServer
	select {}
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
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("Failed to serve: %v", err)
		}
	}()

	//server crashes after 10 seconds :3c
	go func() {
		time.Sleep(time.Second * 10)
		log.Print("!!!Main server has crashed!!!")
		log.Print("Now changing server...")

		//write how far along in time the auction was
		file, err := os.OpenFile("/tmp/logstime.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open file")
		}
		//below 3 lines from chatgpt, write bid as string to log for readability
		_, err = fmt.Fprint(file, time.Since(startTime))
		if err != nil {
			log.Fatal("Failed to write to file:", err)
		}
		os.Exit(0)
	}()
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

		fmt.Fprintf(file, "%s\n", "Highest bid was:")
		//below 3 lines from chatgpt, write bid as string to log for readability
		_, err = fmt.Fprintf(file, "%d\n", HighestBid)
		if err != nil {
			log.Fatal("Failed to write to file:", err)
		}

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

func (server *Server) SendUpdateToFollower(ctx context.Context, request *pb.Request) (*pb.Update, error) {
	if request.AliveMessage {
		return &pb.Update{
			Bid:            HighestBid,
			TimeSinceStart: fmt.Sprint(time.Since(startTime)),
		}, nil
	}
	return &pb.Update{
		Bid:            0,
		TimeSinceStart: "0",
	}, nil
}
