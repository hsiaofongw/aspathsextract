package main

import (
	"fmt"
	"log"

	"sort"

	"example.com/linkanalysis/model"
	"github.com/alecthomas/kong"
)

type NodeScore struct {
	NodeId string
	Score  float64
	Rank   int
}

func calcPageRank() {
	g, err := model.NewGraphFromStdin()
	if err != nil {
		panic(err)
	}

	pr := model.NewPageRank(g, nil)
	for i := 0; i < 100; i++ {
		log.Printf("iteration: %d diff: %f\n", i, pr.UpdateAllNodes())
	}

	nodes := g.GetNodes()
	ns := make([]*NodeScore, 0)
	for _, node := range nodes {
		x := new(NodeScore)
		x.NodeId = node
		x.Score = pr.GetScore(node)
		ns = append(ns, x)
	}

	sort.SliceStable(ns, func(i, j int) bool {
		return ns[i].Score >= ns[j].Score
	})

	fmt.Println("begin page_rank_data:")
	for i := 0; i < len(ns); i++ {
		node := ns[i]
		node.Rank = i
		fmt.Printf("[%d] %s -> %f\n", node.Rank, node.NodeId, node.Score)
	}
}

var CLI struct {
	JSON     bool     `help:"Output in JSON."`
	PageRank struct{} `cmd:"" name:"pagerank" help:"Output PageRank Calculations."`
	Links    struct{} `cmd:"" name:"link" help:"Output Link Analysis."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "rm <path>":
		panic("not implemented")
	case "pagerank":
		calcPageRank()
	default:
		panic(ctx.Command())
	}

	// targetNodeId := ""
	// if len(os.Args) >= 2 && os.Args[1] != "" {
	// 	targetNodeId = os.Args[1]
	// }

	// fmt.Println("begin_data:")

	// fmt.Printf("number_of_nodes: %d\n", g.GetNumNodes())
	// fmt.Printf("number_of_links: %d\n", g.GetNumLinks())

	// if targetNodeId != "" {
	// 	fmt.Printf("target_node: %s\n", targetNodeId)
	// 	inbounds := g.GetInbounds(targetNodeId)
	// 	fmt.Printf("num_inbounds: %d\n", len(inbounds))
	// 	fmt.Println("inbounds:")
	// 	for _, nodeId := range inbounds {
	// 		fmt.Println(nodeId)
	// 	}

	// 	fmt.Printf("num_outbounds: %d\n", g.GetNumOutbounds(targetNodeId))
	// 	outbounds := g.GetOutbounds(targetNodeId)
	// 	fmt.Println("outbounds:")
	// 	for _, nodeId := range outbounds {
	// 		fmt.Println(nodeId)
	// 	}

	// 	fmt.Printf("num_neighbors_total: %d\n", g.GetNumNeighbors(targetNodeId))
	// }

}
