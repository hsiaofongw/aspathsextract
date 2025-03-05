package model

import (
	"math"
)

type PageRankAlgoParams struct {
	// should take 0.85 by default, as suggested by the PageRank paper.
	DumpingFactor float64
}

type PageRank struct {
	g      *Graph
	scores map[string]float64
	params *PageRankAlgoParams
}

func NewPageRank(g *Graph, params *PageRankAlgoParams) *PageRank {
	pr := new(PageRank)
	pr.g = g
	pr.params = params
	if params == nil {
		pr.params = new(PageRankAlgoParams)
		pr.params.DumpingFactor = 0.85
	}
	scores := make(map[string]float64)
	pr.scores = scores
	nodes := g.GetNodes()
	N := len(nodes)
	for _, node := range nodes {
		scores[node] = float64(1) / float64(N)
	}
	return pr
}

func (pr *PageRank) UpdateOneNode(targetNodeId string) float64 {
	inbNodes := pr.g.GetInbounds(targetNodeId)
	prev := pr.scores[targetNodeId]
	var newProb float64 = 0
	N := pr.g.GetNumNodes()
	d := pr.params.DumpingFactor
	for _, inNode := range inbNodes {
		numOuts := pr.g.GetNumOutbounds(inNode)
		prob := pr.scores[inNode]
		if numOuts > 0 {
			contrib := float64(1) / float64(numOuts)
			newProb += contrib * prob * d
		}
	}

	newProb += (1 - d) * (float64(1) / float64(N))
	pr.scores[targetNodeId] = newProb
	diff := prev - newProb
	return diff * diff
}

func (pr *PageRank) UpdateAllNodes() float64 {
	nodes := pr.g.GetNodes()
	var diff float64 = 0
	for _, node := range nodes {
		diff += pr.UpdateOneNode(node)
	}

	return math.Sqrt(diff)
}

func (pr *PageRank) GetScore(node string) float64 {
	return pr.scores[node]
}
