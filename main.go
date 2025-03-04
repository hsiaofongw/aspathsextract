package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
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
	inbounds map[string][]string

	// set of nodeIds
	nodeset map[string]bool
}

func NewGraph() *Graph {
	g := new(Graph)
	g.outlinks = make(map[string]map[string]bool)
	g.inbounds = make(map[string][]string)
	g.nodeset = make(map[string]bool)
	return g
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
	g.inbounds[to] = append(g.inbounds[to], from)
}

// see what nodes are here
func (g *Graph) GetNodes() []string {
	result := make([]string, 0)
	for k := range g.nodeset {
		result = append(result, k)
	}

	return result
}

// see how many nodes are here
func (g *Graph) GetNumNodes() int {
	return len(g.nodeset)
}

// see who is pointing to nodeId
func (g *Graph) GetInbounds(nodeId string) []string {
	if nodes, found := g.inbounds[nodeId]; found {
		return nodes
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

func main() {
	cnt := 0
	lineReader := bufio.NewReader(os.Stdin)

	g := NewGraph()

	targetNodeId := ""
	if len(os.Args) >= 2 && os.Args[1] != "" {
		targetNodeId = os.Args[1]
	}

	for {
		line, err := lineReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) || err == io.EOF {
				log.Println("All links ingested with no error.")
				break
			}
			panic(err)
		}

		line = strings.TrimSpace(line)
		segs := strings.Split(line, ",")
		if len(segs) != 2 {
			log.Printf("Warning: Invalid line: %s\n", line)
			continue
		}

		lhs := segs[0]
		rhs := segs[1]
		log.Printf("[%d]: %s -> %s\n", cnt, lhs, rhs)
		g.AddLink(lhs, rhs)
		cnt++
	}

	fmt.Printf("number_of_nodes: %d\n", g.GetNumNodes())
	allNodes := g.GetNodes()
	fmt.Printf("nodes:\n")
	for _, node := range allNodes {
		fmt.Println(node)
	}
	if targetNodeId != "" {
		fmt.Printf("target_node: %s\n", targetNodeId)
		inbounds := g.GetInbounds(targetNodeId)
		fmt.Printf("num_inbounds: %d\n", len(inbounds))
		fmt.Printf("inbounds:\n")
		for _, nodeId := range inbounds {
			fmt.Println(nodeId)
		}

		fmt.Printf("num_outbounds: %d\n", g.GetNumOutbounds(targetNodeId))
		outbounds := g.GetOutbounds(targetNodeId)
		fmt.Printf("outbounds:\n")
		for _, nodeId := range outbounds {
			fmt.Println(nodeId)
		}
	}
}
