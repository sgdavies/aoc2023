package main

import (
	"fmt"
)

func day21() {
	lines := LinesFromFile("data/21.txt")
	locs := make(map[Point]bool)
	garden := make(map[Point]bool)
	for row, line := range lines {
		for col, r := range line {
			p := Point{row, col}
			if r == 'S' {
				locs[p] = true
				garden[p] = true
			} else if r == '.' {
				garden[p] = true
			}
		}
	}

	for s := 0; s < 64; s++ {
		newLocs := make(map[Point]bool)
		for loc := range locs {
			for _, delta := range []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				np := Point{loc.row + delta.row, loc.col + delta.col}
				if garden[np] {
					newLocs[np] = true
				}
			}
		}
		locs = newLocs
	}

	fmt.Println(len(locs))
}
