package main

import (
	"bufio"
	"fmt"
	"os"
)

type Direction int

const (
	up Direction = iota
	down
	left
	right
)

type beam struct {
	p Point
	d Direction
}

func day16() {
	file, err := os.Open("data/16.txt")
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

	ans := solve16(lines, beam{Point{0, 0}, right})
	fmt.Println(ans)

	// Part two - try all possible starts and return the best
	best := 0
	for c := 0; c < len(lines[0]); c++ {
		ans = solve16(lines, beam{Point{0, c}, down})
		if ans > best {
			best = ans
		}
	}
	for c := 0; c < len(lines[0]); c++ {
		ans = solve16(lines, beam{Point{len(lines) - 1, c}, up})
		if ans > best {
			best = ans
		}
	}
	for r := 0; r < len(lines); r++ {
		ans = solve16(lines, beam{Point{r, 0}, right})
		if ans > best {
			best = ans
		}
	}
	for r := 0; r < len(lines); r++ {
		ans = solve16(lines, beam{Point{r, len(lines[0]) - 1}, left})
		if ans > best {
			best = ans
		}
	}

	fmt.Println(best)
}

func solve16(lines []string, startBeam beam) int {
	grid := make(map[Point]rune)
	for row, line := range lines {
		for col, r := range line {
			grid[Point{row, col}] = r
		}
	}

	energized := make(map[Point]bool)
	beams := []beam{startBeam}
	// Sometimes we get loops - so ignore if we've analyzed any beam before
	seemBeams := make(map[beam]bool)

	for len(beams) > 0 {
		// fmt.Println(len(beams))
		beam := beams[0]
		beams = beams[1:]
		for stillInGrid(beam.p, grid) {
			if seemBeams[beam] {
				break
			} else {
				seemBeams[beam] = true
			}
			energized[beam.p] = true
			// fmt.Println(beam.p)
			if r, present := grid[beam.p]; present && r != '.' {
				newBeam := hitMirror(r, &beam)
				if newBeam != nil {
					beams = append(beams, *newBeam)
				}
			}
			beam.p = beam.p.step(beam.d)
		}
	}
	// fmt.Println(len(energized))
	return len(energized)
}

func stillInGrid(p Point, grid map[Point]rune) bool {
	_, pres := grid[p]
	return pres
}

func (p Point) step(d Direction) Point {
	switch d {
	case up:
		return Point{p.row - 1, p.col}
	case down:
		return Point{p.row + 1, p.col}
	case left:
		return Point{p.row, p.col - 1}
	case right:
		return Point{p.row, p.col + 1}
	default:
		panic("Unexpected direction")
	}
}

func hitMirror(m rune, b *beam) *beam {
	// Modify input beam, and sometimes return an extra beam
	// Pointy end of splitter - no change
	// Flat side of splitter - input changed, plus extra beam returned
	// Angled mirror - input changed
	switch m {
	case '|':
		{
			if b.d == left || b.d == right {
				b.d = up
				return &beam{b.p, down}
			}
		}
	case '-':
		{
			if b.d == up || b.d == down {
				b.d = left
				return &beam{b.p, right}
			}
		}
	case '/':
		{
			if b.d == up {
				b.d = right
			} else if b.d == right {
				b.d = up
			} else if b.d == down {
				b.d = left
			} else if b.d == left {
				b.d = down
			}
		}
	case '\\':
		{
			if b.d == up {
				b.d = left
			} else if b.d == right {
				b.d = down
			} else if b.d == down {
				b.d = right
			} else if b.d == left {
				b.d = up
			}
		}
	}

	return nil
}
