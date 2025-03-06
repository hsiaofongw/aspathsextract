package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"sort"

	"example.com/linkanalysis/model"
	"github.com/alecthomas/kong"
)

type NodeScore struct {
	NodeId       string  `json:"node_id"`
	Score        float64 `json:"score"`
	Rank         int     `json:"rank"`
	NumPeers     int     `json:"num_peers"`
	NumInbounds  int     `json:"num_inbounds"`
	NumOutbounds int     `json:"num_outbounds"`
}

type CommandCtx struct {
	JSON bool
}

type PageRankCmd struct {
	IterationCount *int `arg:"" name:"iters" optional:"" help:"Number of iterations, 100 by default"`
}

func (cmd *PageRankCmd) Run(ctx *CommandCtx) error {
	iters := 100
	if cmd.IterationCount != nil {
		iters = *cmd.IterationCount
	}

	g, err := model.NewGraphFromStdin()
	if err != nil {
		return err
	}

	var output = struct {
		DiffSquarePerIterations []float64    `json:"diff_square_per_iter"`
		NodeAndScores           []*NodeScore `json:"node_and_scores"`
		Timestamp               int64        `json:"generated_at"`
	}{
		DiffSquarePerIterations: make([]float64, 0),
	}

	pr := model.NewPageRank(g, nil)

	iter_idx := 0
	for iters != 0 {
		d := pr.UpdateAllNodes()
		if !ctx.JSON {
			log.Printf("iteration: %d diff: %f\n", iter_idx, d)
		}
		output.DiffSquarePerIterations = append(output.DiffSquarePerIterations, d)

		iters--
		iter_idx++
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

	for i := 0; i < len(ns); i++ {
		ns[i].Rank = i
		ns[i].NumInbounds = len(g.GetInbounds(ns[i].NodeId))
		ns[i].NumOutbounds = g.GetNumOutbounds(ns[i].NodeId)
		ns[i].NumPeers = g.GetNumNeighbors(ns[i].NodeId)
	}

	if !ctx.JSON {
		fmt.Println("begin page_rank_data:")
		for i := 0; i < len(ns); i++ {
			node := ns[i]
			fmt.Printf("[%d] %s -> %f\n", node.Rank, node.NodeId, node.Score)
		}
	}

	output.NodeAndScores = ns

	if ctx.JSON {
		output.Timestamp = time.Now().UnixMilli()
		return json.NewEncoder(os.Stdout).Encode(output)
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

	var output = struct {
		TargetNodeId   string   `json:"target_node"`
		NumInbounds    int      `json:"num_inbounds"`
		NumOutbounds   int      `json:"num_outbounds"`
		Inbounds       []string `json:"inbounds"`
		Outbounds      []string `json:"outbounds"`
		NeighborsTotal int      `json:"num_neighbors_total"`
		Neighbors      []string `json:"neighbors"`
	}{
		TargetNodeId: targetNodeId,
	}

	inbounds := g.GetInbounds(targetNodeId)
	output.NumInbounds = len(inbounds)
	output.Inbounds = inbounds
	output.NumOutbounds = g.GetNumOutbounds(targetNodeId)
	output.Outbounds = g.GetOutbounds(targetNodeId)
	output.NeighborsTotal = g.GetNumNeighbors(targetNodeId)
	output.Neighbors = g.GetNeighbors(targetNodeId)

	return json.NewEncoder(os.Stdout).Encode(output)
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
	JSON     bool        `help:"Forcibly outputs in JSON, otherwise the format of outputs would be command-specific."`
	PageRank PageRankCmd `cmd:"" name:"pagerank" help:"Output PageRank Calculations."`
	Links    LinksCmd    `cmd:"" name:"link" aliases:"links" help:"Output Link Analysis."`
	Overview OverviewCmd `cmd:"" name:"overview" help:"Overviewing the whole graph."`
}

func main() {
	ctx := kong.Parse(&CLI)

	err := ctx.Run(&CommandCtx{JSON: CLI.JSON})
	ctx.FatalIfErrorf(err)
}
