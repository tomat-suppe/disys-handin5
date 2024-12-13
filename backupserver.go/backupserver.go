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
	"google.golang.org/grpc/credentials/insecure"
)

var BidAmount int64
var WinningBidder int32
var HighestBid int64
var startTime time.Time
var bidder *pb.Bidder
var serverCrashed bool = false
var leader bool

type Server struct {
	pb.UnimplementedAuctionServer
}

var TimeAuctionHasRun time.Duration

func main() {
	var server = &Server{}
	leader = false
	TurnOnServer(server)
	for {
		if leader {
			break
		}
		ReceiveHeartBeat()
		time.Sleep(time.Second * 2)
	}

	//keeps server running despite the go-routine .Serve(...) in TurnOnServer
	select {}
}

func ListenForServerCrash() {
	//below 3 lines adapted from https://stackoverflow.com/questions/56336168/golang-check-tcp-port-open
	timeout := time.Second
	_, err := net.DialTimeout("tcp", "localhost:50000", timeout)
	if err != nil {
		serverCrashed = true
	}
	for {
		timeout := time.Second
		_, err := net.DialTimeout("tcp", "localhost:50000", timeout)
		if err != nil {
			serverCrashed = true
		}
		if serverCrashed {
			leader = true
			log.Print("...Leader has shut down, this server is now the Leader...")
			log.Printf("Taking over with highest big: %v and time since auction start: %v", HighestBid, TimeAuctionHasRun.Seconds())
			break
		}
	}

}

// some code from previous hand-ins
func TurnOnServer(server *Server) {

	listener, err := net.Listen("tcp", "localhost:50001")
	if err != nil {
		log.Fatalf("Backupserver failed to listen: %v", err)
	} else {
		log.Printf("Backupserver now listening on: localhost:50001")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuctionServer(grpcServer, server)

	log.Printf("Backupserver is running on : localhost:50001")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	go ListenForServerCrash()

}

// regarding the requirement about bit taking amount as arg: here 'bidder'
// is a standin for this arg, as it has the field 'Bid', with the method .GetBid()
// using this, I can calculate a BidAmount for each bidding, and both server and Client
// keep track of their latest bids.
func (s *Server) Bid(ctx context.Context, in *pb.Bidder) (*pb.BidAccepted, error) {
	if !leader {
		message := "You are speaking to the wrong server.. Try again"
		BidAccepted := &pb.BidAccepted{
			Acceptancemssage: message,
		}

		return BidAccepted, nil
	}
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

		fmt.Fprintf(file, "%s\n", "Highest bid was:")
		//below 3 lines from chatgpt
		_, err = fmt.Fprintf(file, "%d\n", HighestBid)
		if err != nil {
			log.Fatal("Failed to write to file:", err)
		}

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
	if !leader {
		ResultAuctionUpdate := &pb.ResultAuctionUpdate{
			AuctionOverMessage: "Cannot relay auction result, please try again...",
			WinningBid:         -6000,
			WinningBidderId:    -6000,
			//-6000 is an absurd number that could never exist in the auction
			//as go does not support null values for int64s, this is my 'null'
		}

		return ResultAuctionUpdate, nil
	}
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

func ReceiveHeartBeat() {
	conn, err := grpc.NewClient("localhost:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()

	Client := pb.NewAuctionClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	//request a heartbeat with information from leader
	heartbeat, err := Client.SendUpdateToFollower(ctx, &pb.Request{AliveMessage: true})
	if err != nil {

		cancel()
	}
	HighestBid = heartbeat.Bid
	log.Printf("Highest bid at Leader is currently %v, this server is still in Follower state...", heartbeat.Bid)

	TimeAuctionHasRun, _ = time.ParseDuration(heartbeat.TimeSinceStart)
	startTime = time.Now().Add(-TimeAuctionHasRun)
	log.Printf("Auction has run %v seconds, server is still in Follower state...", TimeAuctionHasRun.Seconds())
}
