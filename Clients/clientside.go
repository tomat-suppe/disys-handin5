package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "disys-handin5/proto_files"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var start time.Time
var Client pb.AuctionClient

func main() {
	start = time.Now()
	addr := 50000
	for i := 0; i < 1; i++ {
		addr = addr + i
		addrString := "LocalHost:" + fmt.Sprint(addr)
		bidder := &pb.Bidder{
			BidderId: int32(i),
			Addr:     addrString,
			Bid:      0,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		conn, err := grpc.NewClient(bidder.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		log.Printf("Connected bidder #%v!", bidder.BidderId)
		Client = pb.NewAuctionClient(conn)
		//TurnOnClient(bidder)

		for i := 0; i < 50; i++ {
			go SendBid(ctx, bidder)
			b, err := Client.Bid(ctx, bidder)
			if err != nil {
				log.Fatalf("Failed to call bid! %v", err)
			}
			log.Printf(b.GetAcceptancemssage())
			r, err := Client.Result(ctx, bidder)
			if err != nil {
				log.Fatalf("Failed to call bid! %v", err)
			}
			log.Printf(r.GetAuctionOverMessage())

		}
		defer cancel()
		defer conn.Close()
	}
}

func SendBid(ctx context.Context, bidder *pb.Bidder) (bidder1 *pb.Bidder, err error) {
	return bidder, nil
}

/*func TurnOnClient(bidder *pb.Bidder) {

}*/
/*
func Bid(ctx context.Context, bidder *pb.Bidder) (*pb.BidAccepted, error) {

}

func Result(ctx context.Context, bidder *pb.Bidder) (*pb.ResultAuctionUpdate, error) {

}
*/
