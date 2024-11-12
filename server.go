package main

import(
	"bufio"
)

type Node struct {
	NodeID int
	addr string

}

var Scanner bufio.Scanner
var HighestBid int = 0
var FinalBid int = 0
var TimeSteps int = 0 //maybe this needs to be an actual time component

func main(){
	nodeaddr := []string{
		"localhost:50051", // Node 1
		"localhost:50052", // Node 2
		"localhost:50053", // Node 3
	}

	for i := 0; i < 3 ; i++{
		node := Node{
			NodeID: i,
			addr: nodeaddr[i]
		}
		node.Serve()
	}
	//time start
	for{
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

func (Node *Node) bid(amount int) (ack string){
	//receive a bid
	//update either final bid or highest bid
	//depending on time passed.

	//return either "fail", "success" or "exception"
	//success when bid is bigger than Highest bid
	//fail when bid is equal to or less than highest bid

	//exception when bid has timeouted?
}

func (Node *Node) result() (outcome string){
	//if TimeSteps => arbitrary time (100?) {
	//return log.Printf("Auction has ended with final bid being %s. Winner was bidder number %v", HighestBid, Node.NodeID)
	//} else {
	//return log.Printf("Current highest bid is %s by bidder number %v", CurrentBid, Node.NodeID)
	//}
}

func (Node *Node) Serve() {
	/* from previous hand-in:
	listener, err := net.Listen("tcp", node.listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	} else {
		log.Printf("Node %v now listening on: %s", node.NodeID, node.listenAddr)
	}
	node.server = grpc.NewServer()
	//server := Server{}
	pb.RegisterMutualExclusionServer(node.server, node)

	log.Printf("Node %v is running on : %s ...", node.NodeID, node.Addr)
	if err := node.server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}*/
}