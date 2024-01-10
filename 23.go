package main

import "fmt"

func day23() {
	lines := LinesFromFile("data/23.txt")
	grid := make(map[Point]rune)
	var start, end Point
	for row, line := range lines {
		for col, r := range line {
			if r != '#' {
				p := Point{row, col}
				grid[p] = r

				if row == 0 {
					start = p
				} else {
					end = p // continually overridden
				}
			}
		}
	}

	fmt.Println(recurse23(start, 0, map[Point]bool{start: true}, grid, end, true))

	// Longest path problem is NP-hard, and on the large grid this takes too long:
	// fmt.Println(recurse23(start, 0, map[Point]bool{start: true}, grid, end, false))

	// Simplify things by reducing to nodes (intersections) and edges (lengths between
	// intersections) so the size of the network is small enough to solve quickly.
	graph := make(map[Point]map[Point]int)

	toCheck := []exploreStep{{start, start, Point{start.row + 1, start.col}, 1}}
	seen := map[Point]bool{start: true}
	for len(toCheck) > 0 {
		tc := toCheck[0]
		toCheck = toCheck[1:]
		node := tc.curTile
		nextNs := []Point{}
		for _, n := range neighbours23(node, grid[node], false) {
			if _, pres := grid[n]; pres && n != tc.lastTile {
				nextNs = append(nextNs, n)
			}
		}

		if len(nextNs) == 0 {
			// Dead end (should only happen at the end)
			// fmt.Println("Dead end: ", node)
			if node == end {
				addEdge(graph, tc.start, node, tc.len)
				addEdge(graph, node, tc.start, tc.len)
			} else {
				panic("I wasn't expecting a dead end")
			}
		} else if len(nextNs) > 1 {
			addEdge(graph, tc.start, node, tc.len)
			addEdge(graph, node, tc.start, tc.len)

			// seed starting points from the next node - assuming we've
			// not seeded from there already
			if !seen[node] {
				seen[node] = true
				for _, nextN := range nextNs {
					toCheck = append(toCheck, exploreStep{node, node, nextN, 1})
				}
			}
		} else {
			toCheck = append(toCheck, exploreStep{tc.start, node, nextNs[0], tc.len + 1})
		}
	}

	// fmt.Println(len(graph))

	// Now solve from here
	// fmt.Println(r23b(start, 0, map[Point]bool{start: true}, graph, end))

	// ... but to speed up further, solve to the last node before the end, then add
	// on the final distance.  Because once you've got to the penultimate node, you
	// have to go to the end - you can't come back onto the node.  This saves
	// exploring down some dead ends, for a ~2x speedup.
	if len(graph[end]) != 1 {
		panic("More than one node connected to end!?")
	}
	var penultimate Point
	var penLength int
	for k, v := range graph[end] {
		penultimate = k
		penLength = v
	}
	fmt.Println(r23b(start, 0, map[Point]bool{start: true}, graph, penultimate) + penLength)
}

func addEdge(graph map[Point]map[Point]int, a Point, b Point, length int) {
	if _, pres := graph[a]; !pres {
		graph[a] = make(map[Point]int)
	}
	graph[a][b] = length
}

type edge struct {
	otherNode Point
	len       int
}
type exploreStep struct {
	start    Point
	lastTile Point
	curTile  Point
	len      int
}

func r23b(p Point, len int, seen map[Point]bool, graph map[Point]map[Point]int, end Point) int {
	if p == end {
		return len
	}
	seen[p] = true
	best := 0
	for np, l := range graph[p] {
		if !seen[np] {
			ans := r23b(np, len+l, seen, graph, end)
			if ans > best {
				best = ans
			}
		}
	}
	seen[p] = false
	return best
}

func recurse23(p Point, len int, seen map[Point]bool, grid map[Point]rune, end Point, slippery bool) int {
	if p == end {
		// fmt.Println("\t", len) // For checking the example only!
		return len
	}
	seen[p] = true

	best := 0
	newLen := len + 1

	for _, np := range neighbours23(p, grid[p], slippery) {
		if _, pres := grid[np]; pres && !seen[np] {
			ans := recurse23(np, newLen, seen, grid, end, slippery)
			if ans > best {
				best = ans
			}
		}
	}

	seen[p] = false
	return best
}

func neighbours23(p Point, r rune, slippery bool) []Point {
	if slippery {
		switch r {
		case '.':
			{
				return []Point{{p.row - 1, p.col}, {p.row + 1, p.col}, {p.row, p.col - 1}, {p.row, p.col + 1}}
			}
		case '^':
			{
				return []Point{{p.row - 1, p.col}}
			}
		case 'v':
			{
				return []Point{{p.row + 1, p.col}}
			}
		case '>':
			{
				return []Point{{p.row, p.col + 1}}
			}
		case '<':
			{
				return []Point{{p.row, p.col - 1}}
			}
		default:
			{
				panic("Unexpected rune")
			}
		}
	}

	return []Point{{p.row - 1, p.col}, {p.row + 1, p.col}, {p.row, p.col - 1}, {p.row, p.col + 1}}
}
