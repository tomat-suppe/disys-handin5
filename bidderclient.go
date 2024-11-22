package main

import (
	"bufio"
	"log"
)

var Scanner bufio.Scanner
var HighestBid int = 0
var FinalBid int = 0

type Node struct {
	NodeID int
	addr   string
}

func main() {
	nodeaddr := []string{
		"localhost:50051", // Node 1
		"localhost:50052", // Node 2
		"localhost:50053", // Node 3
	}

	/*	for i := 0; i < 3 ; i++{
		node := Node{
			NodeID: i,
			addr: nodeaddr[i]
		}
		node.Serve()
	}*/
	//instead of above, simply run on 1 Node, then if this Node fails (kill it at a certain
	//point, maybe time 50%), it should change over to Node 2.
	//Then Node2 highestbid = whatever was logged from Node1
	//time start

	for {
		input := Scanner.Text()
		if input == "bid" {
			log.Println("Enter amount to bid:")
			amount := Scanner.Text()
			Node.bid(amount)
		} else if input == "result" {
			Node.result()
		}
	}
}

func (Node *Node) bid(amount int) (ack string) {
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

func (Node *Node) result() (outcome string) {
	//if TimeSteps => arbitrary time (100?) {
	//return log.Printf("Auction has ended with final bid being %s. Winner was bidder number %v", HighestBid, Node.NodeID)
	//} else {
	//return log.Printf("Current highest bid is %s by bidder number %v", CurrentBid, Node.NodeID)
	//}
}

/*func Bid (ctx context.Context, bidder *pb.Bidder) (*pb.AuctionUpdate, error){
}*/
