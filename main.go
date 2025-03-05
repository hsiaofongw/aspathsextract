package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"sort"

	"example.com/linkanalysis/model"
	"github.com/alecthomas/kong"
)

type NodeScore struct {
	NodeId string
	Score  float64
	Rank   int
}

type CommandCtx struct {
	JSON bool
}

type PageRankCmd struct{}

func (cmd *PageRankCmd) Run(ctx *CommandCtx) error {
	g, err := model.NewGraphFromStdin()
	if err != nil {
		return err
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

	return nil
}

type LinksCmd struct {
	TargetNodeId string `arg:"" name:"targetnodeid" help:"Id of Target Node to analyze"`
}

func (cmd *LinksCmd) Run(ctx *CommandCtx) error {
	targetNodeId := cmd.TargetNodeId
	if targetNodeId == "" {
		return errors.New("empty node id")
	}

	g, err := model.NewGraphFromStdin()
	if err != nil {
		return err
	}

	fmt.Println("begin_data:")
	fmt.Printf("target_node: %s\n", targetNodeId)
	inbounds := g.GetInbounds(targetNodeId)
	fmt.Printf("num_inbounds: %d\n", len(inbounds))
	fmt.Println("inbounds:")
	for _, nodeId := range inbounds {
		fmt.Println(nodeId)
	}

	fmt.Printf("num_outbounds: %d\n", g.GetNumOutbounds(targetNodeId))
	outbounds := g.GetOutbounds(targetNodeId)
	fmt.Println("outbounds:")
	for _, nodeId := range outbounds {
		fmt.Println(nodeId)
	}

	fmt.Printf("num_neighbors_total: %d\n", g.GetNumNeighbors(targetNodeId))

	return nil
}

type OverviewCmd struct{}

func (cmd *OverviewCmd) Run(ctx *CommandCtx) error {
	g, err := model.NewGraphFromStdin()
	if err != nil {
		return err
	}

	var output = struct {
		NumNodes int `json:"number_of_nodes"`
		NumLinks int `json:"number_of_links"`
	}{NumNodes: g.GetNumNodes(), NumLinks: g.GetNumLinks()}

	return json.NewEncoder(os.Stdout).Encode(output)
}

var CLI struct {
	JSON     bool        `help:"Output in JSON."`
	PageRank PageRankCmd `cmd:"" name:"pagerank" help:"Output PageRank Calculations."`
	Links    LinksCmd    `cmd:"" name:"link" help:"Output Link Analysis."`
	Overview OverviewCmd `cmd:"" name:"overview" help:"Overviewing the whole graph."`
}

func main() {
	ctx := kong.Parse(&CLI)

	err := ctx.Run(&CommandCtx{JSON: CLI.JSON})
	ctx.FatalIfErrorf(err)
}
