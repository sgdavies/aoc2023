package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// This file contains 4 different implementations of various speeds
// All are based on the Karger probabilistic algorithm for finding the minimum cut.
//
// Table shows iteration count and total time for 5 repeated solves, then mean stats:
// 					karger	karger2	karger3	kragerUF
// 	iterations		 513	 135	 624	  59
// 	per-run			 324	 295	 469	  80
// 					 347	 517	  68	 357
// 					 380	1101	 176	 215
// 					  21	 644	 377	 594
// total iter		1585	2692	1714	1305
// mean iter		 317	 538.4	 342.8	 261
// seconds (5 runs)	 340	 425	  24	   4
// iter/sec			   4.7	   6.3	  71.4	 326
// msec/iter		 214	 158	  14	   3
// expected solve time (seconds)
// 					  61	  83	   5.4	   1.2
//
// karger, karger2 and karger3 are broadly similar.
// Each picks an edge at random, creates a new node to represent the joined nodes,
// then iterates through the edges list and updates all edges that referred to the old nodes
// to refer to the new one. Thus the time complexity is roughly O(N*E) - we remove 1 node each
// time round the loop until only 2 remain, and on each time round we look at all the edges.
// We achieve a >10x speedup by:
// karger - naive implementation
// karger2 - maintain a single edge list and flag edges as dead (rather than recreating list) (25% boost)
// karger3 - use int ids for all nodes rather than strings (further 90% boost)
//
// kargerUF uses the Union-Find data structure.  This results in a considerable speed-up - now
// we loop the same number of times, but we only need to do an O(1) operation on a single pair of
// nodes.  So net time is O(N).
// This impl is the "naive" version, using strings - we'd expect a further boost by switching to ints.

func day25() {
	nodes := make(map[string]bool)
	edges := [][2]string{} // (nodeA, nodeB)
	nodes3 := make(map[int]bool)
	node2id3 := make(map[string]int)
	edges3 := [][2]int{} // (nodeA, nodeB)
	nextNodeI := 1
	for _, line := range LinesFromFile("data/25.txt") {
		parts := strings.Split(line, ": ")
		nodeA := parts[0]
		if _, ok := node2id3[nodeA]; !ok {
			node2id3[nodeA] = nextNodeI
			nextNodeI += 1
		}
		nodes[nodeA] = true
		nodes3[node2id3[nodeA]] = true
		for _, nodeB := range strings.Split(parts[1], " ") {
			nodeA := parts[0]
			if _, ok := node2id3[nodeB]; !ok {
				node2id3[nodeB] = nextNodeI
				nextNodeI += 1
			}
			nodes[nodeB] = true
			nodes3[node2id3[nodeB]] = true
			edges = append(edges, [2]string{nodeA, nodeB})
			edges3 = append(edges3, [2]int{node2id3[nodeA], node2id3[nodeB]})
		}
	}
	fmt.Println(len(nodes), len(edges))

	for iterations := 1; true; iterations++ {
		remainingNodes, remainingEdges := kragerUF(nodes, &edges)
		// remainingNodes, remainingEdges := krager3(nodes3, &edges3, nextNodeI)
		if remainingEdges == 3 {
			if len(remainingNodes) != 2 {
				fmt.Println(remainingNodes)
				panic("Should only have 2 nodes!")
			}
			ans := 1
			for _, v := range remainingNodes {
				ans *= v
			}
			fmt.Println(ans, "(took", iterations, "iterations)")
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

func krager2(constNodeNames map[string]bool, constEdges *[][2]string) (map[string]int, int) {
	nodes := make(map[string]int) // Node ID : number of conglomorated nodes
	for k := range constNodeNames {
		nodes[k] = 1
	}
	edges := make([]*[3]string, len(*constEdges))
	for i, e := range *constEdges {
		edges[i] = &[3]string{e[0], e[1], "L"} // L=live (vs D=dead)
	}
	rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
	next := 0

	for len(nodes) > 2 {
		var re [3]string // removed edge
		for {
			re = *edges[next]
			next += 1
			if re[2] == "L" {
				break
			}
		}
		newNodeName := re[0] + re[1]
		// Join the two nodes
		nodes[newNodeName] = nodes[re[0]] + nodes[re[1]]
		delete(nodes, re[0])
		delete(nodes, re[1])
		for i, oe := range edges[next:] { // oe = old edge
			ei := i + next
			if (re[0] == oe[0] && re[1] == oe[1]) || (re[0] == oe[1] && re[1] == oe[0]) {
				// we're removing this edge
				edges[ei][2] = "D"
			}

			if re[0] == oe[0] || re[1] == oe[0] {
				// 0th element needs to be replaced
				edges[ei][0] = newNodeName
			} else if re[0] == oe[1] || re[1] == oe[1] {
				// 1st element needs to be replaced
				edges[ei][1] = newNodeName
			}
		}
	}

	// fmt.Println("Found a cut:", len(edges))
	ans := 0
	for _, e := range edges[next:] {
		if e[2] == "L" {
			ans += 1
		}
	}
	return nodes, ans
}

func krager3(constNodeNames map[int]bool, constEdges *[][2]int, nextName int) (map[int]int, int) {
	nodes := make(map[int]int) // Node ID : number of conglomorated nodes
	for k := range constNodeNames {
		nodes[k] = 1
	}
	edges := make([]*[3]int, len(*constEdges))
	for i, e := range *constEdges {
		edges[i] = &[3]int{e[0], e[1], 1} // 1=live (vs 0=dead)
	}
	rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
	next := 0

	for len(nodes) > 2 {
		var re [3]int // removed edge
		for {
			re = *edges[next]
			next += 1
			if re[2] == 1 {
				break
			}
		}
		newNodeName := nextName
		nextName += 1
		// Join the two nodes
		nodes[newNodeName] = nodes[re[0]] + nodes[re[1]]
		delete(nodes, re[0])
		delete(nodes, re[1])
		for i, oe := range edges[next:] { // oe = old edge
			ei := i + next
			if (re[0] == oe[0] && re[1] == oe[1]) || (re[0] == oe[1] && re[1] == oe[0]) {
				// we're removing this edge
				edges[ei][2] = 0
			}

			if re[0] == oe[0] || re[1] == oe[0] {
				// 0th element needs to be replaced
				edges[ei][0] = newNodeName
			} else if re[0] == oe[1] || re[1] == oe[1] {
				// 1st element needs to be replaced
				edges[ei][1] = newNodeName
			}
		}
	}

	// fmt.Println("Found a cut:", len(edges))
	ans := 0
	for _, e := range edges[next:] {
		if e[2] == 1 {
			ans += 1
		}
	}
	return nodes, ans
}

func kragerUF(constNodeNames map[string]bool, constEdges *[][2]string) (map[string]int, int) {
	// Use union-find datastructure to speed things up
	parents := make(map[string]string)
	sizes := make(map[string]int)
	for k := range constNodeNames {
		parents[k] = k
		sizes[k] = 1
	}

	edges := make([]*[2]string, len(*constEdges))
	for i, e := range *constEdges {
		edges[i] = &[2]string{e[0], e[1]}
	}
	rand.Shuffle(len(edges), func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })

	for _, edge := range edges {
		ufUnion(edge, &parents, &sizes)

		if len(sizes) == 2 {
			break
		}
	}

	minCut := 0
	for _, e := range edges {
		if ufFind(e[0], &parents) != ufFind(e[1], &parents) {
			minCut += 1
		}
	}
	return sizes, minCut
}

func ufUnion(edge *[2]string, parents *map[string]string, sizes *map[string]int) {
	a := ufFind(edge[0], parents)
	b := ufFind(edge[1], parents)

	if a == b {
		return // Already in same group
	}

	// Join together. Join smaller group onto larger.
	var big, small string
	if (*sizes)[a] < (*sizes)[b] {
		big, small = b, a
	} else {
		big, small = a, b
	}
	(*parents)[small] = big
	(*sizes)[big] += (*sizes)[small]
	delete(*sizes, small)
}

func ufFind(node string, parents *map[string]string) string {
	// find
	pCand := node
	n := (*parents)[pCand]
	for pCand != n {
		pCand = n
		n = (*parents)[pCand]
	}
	root := n

	// path compression
	pCand = node
	for pCand != n {
		pCand = n
		n = (*parents)[pCand]
		(*parents)[pCand] = root
	}

	return root
}
