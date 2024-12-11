package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/tomat-suppe/disys-handin5/proto_files"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var start time.Time
var ListOfBidders []*pb.Bidder

func main() {
	start = time.Now()
	addr1 := "localhost:" + fmt.Sprint(50000)
	addr2 := "localhost:" + fmt.Sprint(50001)
	var addresses []string
	addresses = append(addresses, addr1)
	addresses = append(addresses, addr2)

	for i := 0; i < 3; i++ {
		bidder := &pb.Bidder{
			BidderId: int32(i),
			Addr:     addr1,
			Bid:      rand.Int63n(10),
		}
		ListOfBidders = append(ListOfBidders, bidder)
	}
	for {
		for _, bidder := range ListOfBidders {
			_, err := net.DialTimeout("tcp", addr1, time.Second)
			if err != nil {
				log.Print("test")
				addresses = []string{addr2}
			}
			bidder.Addr = addresses[rand.Intn(len(addresses))]
			conn, err := grpc.NewClient(bidder.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Did not connect: %v", err)
			}

			defer conn.Close()

			log.Printf("Bidder #%v is acting!", bidder.BidderId)
			Client := pb.NewAuctionClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			b, err := Client.Bid(ctx, bidder)
			if err != nil {
				log.Fatalf("Failed to call bid! %v", err)
			}
			bidder.Bid = bidder.Bid + rand.Int63n(1000)
			log.Printf("Bidder #%v: %s", bidder.BidderId, b.GetAcceptancemssage())
			r, err := Client.Result(ctx, bidder)
			if err != nil {
				log.Fatalf("Failed to call result! %v", err)
			}
			log.Printf(r.GetAuctionOverMessage())

			time.Sleep(time.Second)
		}
	}
}
