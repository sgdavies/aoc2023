package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func day25() {
	nodes := make(map[string]bool)
	edges := [][2]string{} // (nodeA, nodeB)
	for _, line := range LinesFromFile("data/25.txt") {
		parts := strings.Split(line, ": ")
		nodeA := parts[0]
		nodes[nodeA] = true
		for _, nodeB := range strings.Split(parts[1], " ") {
			nodes[nodeB] = true
			edges = append(edges, [2]string{nodeA, nodeB})
		}
	}

	for {
		remainingNodes, remainingEdges := krager(nodes, &edges)
		if remainingEdges == 3 {
			if len(remainingNodes) != 2 {
				fmt.Println(remainingNodes)
				panic("Should only have 2 nodes!")
			}
			ans := 1
			for _, v := range remainingNodes {
				ans *= v
			}
			fmt.Println(ans)
			break
		}
	}
}

func krager(constNodeNames map[string]bool, constEdges *[][2]string) (map[string]int, int) {
	nodes := make(map[string]int) // Node ID : number of conglomorated nodes
	for k := range constNodeNames {
		nodes[k] = 1
	}
	edges := *constEdges

	for len(nodes) > 2 {
		r := rand.Intn(len(edges))
		re := edges[r] // removed edge
		newNodeName := re[0] + re[1]
		// Join the two nodes
		nodes[newNodeName] = nodes[re[0]] + nodes[re[1]]
		delete(nodes, re[0])
		delete(nodes, re[1])
		// Remove the old edges, and swap in the new node in place of either old node on any other edges
		newEdges := make([][2]string, len(edges))
		index := 0
		for _, oe := range edges { // oe = old edge
			if (re[0] == oe[0] && re[1] == oe[1]) || (re[0] == oe[1] && re[1] == oe[0]) {
				// we're removing this edge
				continue
			}

			if re[0] == oe[0] || re[1] == oe[0] {
				// 0th element needs to be replaced
				newEdges[index] = [2]string{newNodeName, oe[1]}
			} else if re[0] == oe[1] || re[1] == oe[1] {
				// 1st element needs to be replaced
				newEdges[index] = [2]string{oe[0], newNodeName}
			} else {
				newEdges[index] = oe
			}
			index += 1
		}
		edges = newEdges[:index]
	}

	// fmt.Println("Found a cut:", len(edges))
	return nodes, len(edges)
}
