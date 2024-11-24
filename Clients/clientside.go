package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/tomat-suppe/disys-handin5/proto_files"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var start time.Time
var ListOfBidders []*pb.Bidder

//var Client pb.AuctionClient

func main() {
	start = time.Now()
	addr := 50000
	for i := 0; i < 3; i++ {
		//addr = addr + i
		addrString := "localhost:" + fmt.Sprint(addr)
		bidder := &pb.Bidder{
			BidderId: int32(i),
			Addr:     addrString,
			Bid:      0,
		}
		ListOfBidders = append(ListOfBidders, bidder)
	}
	for i := 0; i < 50; i++ {
		for _, bidder := range ListOfBidders {
			conn, err := grpc.NewClient(bidder.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Did not connect: %v", err)
			}

			defer conn.Close()
			log.Printf("Bidder #%v is acting!", bidder.BidderId)
			Client := pb.NewAuctionClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			//TurnOnClient(bidder)

			go SendBid(ctx, bidder)
			b, err := Client.Bid(ctx, bidder)
			if err != nil {
				log.Fatalf("Failed to call bid! %v", err)
			}
			log.Printf("Bidder #%v: %s", bidder.BidderId, b.GetAcceptancemssage())
			r, err := Client.Result(ctx, bidder)
			if err != nil {
				log.Fatalf("Failed to call result! %v", err)
			}
			log.Printf(r.GetAuctionOverMessage())

		}

	}
}

func SendBid(ctx context.Context, bidder *pb.Bidder) (bidder1 *pb.Bidder, err error) {
	return bidder, nil
}
