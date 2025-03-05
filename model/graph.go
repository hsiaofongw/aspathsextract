package model

import (
	"bufio"
	"errors"
	"io"

	"os"
	"strings"
)

type Graph struct {
	// nodeId -> outlinks of nodeId
	// e.g.: 'a' -> 'b', 'a' -> 'c', 'a' -> 'd' then
	// map['a'] = { 'b', 'c', 'd' }
	outlinks map[string]map[string]bool

	// nodeId -> set of nodes that points to it
	// e.g.: map['a'] = { 'b', 'c', 'd' } iif. 'b', 'c', 'd' all has links points to 'a'.
	inbounds map[string]map[string]bool

	// set of nodeIds
	nodeset map[string]bool

	// set of any inbound(or outbound) neighbors
	neighbors map[string]map[string]bool

	// num links
	numlinks int
}

func NewGraph() *Graph {
	g := new(Graph)
	g.outlinks = make(map[string]map[string]bool)
	g.inbounds = make(map[string]map[string]bool)
	g.neighbors = make(map[string]map[string]bool)
	g.nodeset = make(map[string]bool)
	g.numlinks = 0
	return g
}

func NewGraphFromStdin() (*Graph, error) {
	lineReader := bufio.NewReader(os.Stdin)
	g := NewGraph()
	for {
		line, err := lineReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) || err == io.EOF {
				break
			}
			return nil, err
		}

		line = strings.TrimSpace(line)
		segs := strings.Split(line, ",")
		for i := 0; i < len(segs)-1; i++ {
			lhs := segs[i]
			rhs := segs[i+1]
			g.AddLink(lhs, rhs)
		}
	}

	return g, nil
}

func (g *Graph) AddLink(from, to string) {
	g.nodeset[from] = true
	g.nodeset[to] = true

	if submap, found := g.outlinks[from]; found {
		if _, found := submap[to]; found {
			return
		}
	} else {
		g.outlinks[from] = make(map[string]bool)
	}

	g.outlinks[from][to] = true
	if _, found := g.inbounds[to]; !found {
		g.inbounds[to] = make(map[string]bool)
	}
	g.inbounds[to][from] = true
	if _, found := g.neighbors[from]; !found {
		g.neighbors[from] = make(map[string]bool)
	}
	g.neighbors[from][to] = true
	if _, found := g.neighbors[to]; !found {
		g.neighbors[to] = make(map[string]bool)
	}
	g.neighbors[to][from] = true
	g.numlinks++
}

// see what nodes are here
func (g *Graph) GetNodes() []string {
	result := make([]string, 0)
	for k := range g.nodeset {
		result = append(result, k)
	}

	return result
}

// see how many links
func (g *Graph) GetNumLinks() int {
	return g.numlinks
}

// see how many nodes are here
func (g *Graph) GetNumNodes() int {
	return len(g.nodeset)
}

// see who is pointing to nodeId
func (g *Graph) GetInbounds(nodeId string) []string {
	if nodes, found := g.inbounds[nodeId]; found {
		res := make([]string, 0)
		for n := range nodes {
			res = append(res, n)
		}
		return res
	}

	result := make([]string, 0)
	return result
}

// see who nodeId is pointing to
func (g *Graph) GetOutbounds(nodeId string) []string {
	result := make([]string, 0)
	if submap, found := g.outlinks[nodeId]; found {
		for k := range submap {
			result = append(result, k)
		}
	}

	return result
}

// count the number of outlinks, useful at calculating PageRank
func (g *Graph) GetNumOutbounds(nodeId string) int {
	if submap, found := g.outlinks[nodeId]; found {
		return len(submap)
	}
	return 0
}

// count neighbors, regardless of inbound or outbound
func (g *Graph) GetNumNeighbors(nodeId string) int {
	if nodes, found := g.neighbors[nodeId]; found {
		return len(nodes)
	}

	return 0
}

// list all neighbors, regardless of inbound or outbound
func (g *Graph) GetNeighbors(nodeId string) []string {
	res := make([]string, 0)
	if nodes, found := g.neighbors[nodeId]; found {
		for n := range nodes {
			res = append(res, n)
		}
	}

	return res
}
