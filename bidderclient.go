package main

import (
	"bufio"
	"log"
	"time"
)

type HighestBid struct {
	Bid int64
	BidderId int32
}

var start time
/*type Node struct {
	NodeID int
	addr   string
	Bid int64
}*/

func main() {
	/*nodeaddr := []string{
		"localhost:50051", // Node 1
		"localhost:50052", // Node 2
		"localhost:50053", // Node 3
	}*/
		start := time.Now()
		for i := 0; i < 3 ; i++{
		bidder := Bidder{
			NodeID: i,
			addr: "localhost:50051",
			Bid: 0,
		}
		
		bidder.TurnOnClient()

		go bidder.bid()
		go bidder.result()
	}
}

func (node *Node) TurnOnClient() {
	conn, err := grpc.NewClient(node.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
		//node.client = pb.NewMutualExclusionClient(conn)
}

func (Node *Node) Bid (ctx context.Context, bidder *pb.Bidder) (stream *pb.AuctionUpdate, error) {
	Node.bid +5

	if Node.bid > HighestBid.Bid {
		log.Printf("%v has updated highest bid with bid: %s", Node.NodeID, Node.bid)
		return AuctionUpdate {
			HighestBid: Node.bid;
    		HighestBidderId: Node.NodeID;
		}
	} else {
		log.Printf("%v tried to bid %s but bid was too low!", Node.NodeID, Node.bid)
	}

	//receive a bid
	//update either final bid or highest bid
	//depending on time passed.

	//return either "fail", "success" or "exception"
	//success when bid is bigger than Highest bid
	//fail when bid is equal to or less than highest bid

	//exception when bid has timeouted?

	//if time < 100 {
	//	if amount > HighestBid {
	//		HighestBid = amount
	//	Maybe reformat all this shit to only return the three required words
	// 	Then handle logic in main or a formatter class?
	//		log.Printf("Success: Bid from bidder number %v has been accepted, new highest bid is %s.", Node.NodeID, HighestBid)
	//	} else{
	//		log.Printf("Failure: ")
	//	}
	//} else {
	// return log.Printf("Exception: Auction has ended. Final bid was %s", FinalBid)
	//}
}

func (Node *Node) Result (ctx context.Context, bidder *pb.Bidder) (stream *pb.ResultAuctionUpdate, error) {
	
	if time.Since(start) > 500 {
		log.Printf("Auction has ended with final bid being %s. Winner was bidder number %v", HighestBid.Bid, HighestBid.BidderId)
		return ResultAuctionUpdate{
			AuctionOver: true,
			WinningBid: HighestBid.Bid,
			WinningBidderId: HighestBid.BidderId,
		}
	} else {
		log.Printf("Current highest bid is %s by bidder number %v", CurrentBid, Node.NodeID)
		return ResultAuctionUpdate{
			AuctionOver: false,
			WinningBid: HighestBid.Bid,
			WinningBidderId: HighestBid.BidderId,
		}
	}
}

/*func Bid (ctx context.Context, bidder *pb.Bidder) (*pb.AuctionUpdate, error){
}*/
