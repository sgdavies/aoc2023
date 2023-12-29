package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Path struct {
	p     Point
	rl    int // run length (must be <= 3)
	lastP Point
}

func day17() {
	// 	testLines := LinesFromFile("data/17.ex")
	// 	solve17(testLines, false)
	// 	solve17(testLines, true)
	// 	solve17(strings.Split(`111111111111
	// 999999999991
	// 999999999991
	// 999999999991
	// 999999999991`, "\n"), true)
	lines := LinesFromFile("data/17.txt")
	solve17(lines, false)
	solve17(lines, true)
}

func solve17(lines []string, ultraCrucibles bool) {
	maxLine := 3
	if ultraCrucibles {
		maxLine = 10
	}

	grid := make(map[Point]int) // cost of each grid square
	for row, line := range lines {
		for col, r := range line {
			cost, _ := strconv.Atoi(string(r))
			grid[Point{row, col}] = cost
		}
	}
	// fmt.Println("Grid size: ", len(lines), len(lines[0]))

	// Djikstra - but each node in graph is (point, last point, distance)
	startPath := Path{Point{0, 0}, 0, Point{-1, -1}} // lastP not on same row or col
	end := Point{len(lines) - 1, len(lines[0]) - 1}

	best := map[Path]int{startPath: 0}
	costs := make(map[int]map[Path]bool)
	costs[0] = map[Path]bool{startPath: true}

	for {
		keys := make([]int, len(costs))
		i := 0
		for k := range costs {
			keys[i] = k
			i++
		}
		slices.Sort(keys)
		curCost := keys[0]

		var curPath Path
		for k := range costs[curCost] {
			curPath = k
			break
		}

		if curPath.p == end {
			// With ultra-crucibles, can't stop until we've gone at least 4
			if !(ultraCrucibles && curPath.rl < 4) {
				fmt.Println(curCost)
				break
			}
		}

		delete(costs[curCost], curPath)
		if len(costs[curCost]) == 0 {
			delete(costs, curCost)
		}

		for _, delta := range []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			newP := Point{curPath.p.row + delta.row, curPath.p.col + delta.col}
			if newP == curPath.lastP {
				// Can only go straight or turn L or R - no U-turns allowed
				continue
			}

			if ultraCrucibles && curPath.rl < 4 {
				if !(curPath.lastP.row < 0 || curPath.lastP.row == newP.row || curPath.lastP.col == newP.col) {
					// Must go at least 4 in a row
					// (plus hack to cope with first step)
					continue
				}
			}

			if cost, present := grid[newP]; present {
				var newRL int
				if curPath.lastP.row == newP.row || curPath.lastP.col == newP.col {
					newRL = curPath.rl + 1
				} else {
					newRL = 1
				}

				if newRL > maxLine {
					continue
				}

				newPath := Path{newP, newRL, curPath.p}
				newCost := curCost + cost

				curBest, present := best[newPath] // defaults to 0
				if !present || newCost < curBest {
					best[newPath] = newCost
					delete(costs[curBest], newPath) // nop if not present
					if len(costs[curCost]) == 0 {
						delete(costs, curCost)
					}
					if _, pres := costs[newCost]; !pres {
						costs[newCost] = make(map[Path]bool)
					}
					costs[newCost][newPath] = true
				}
			}
		}
	}
}

func LinesFromFile(fname string) []string {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
