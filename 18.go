package main

import (
	"fmt"
	"strconv"
	"strings"
)

func day18() {
	lines := LinesFromFile("data/18.txt")

	solve18(lines, false)
	solve18(lines, true)
}

func parse18(line string, partTwo bool) (byte, int) {
	parts := strings.Split(line, " ")
	if !partTwo {
		v, _ := strconv.Atoi(parts[1])
		return parts[0][0], v
	} else {
		// (#70c710) -> (#, 70c71->461937, 0->R, )
		hex := parts[2]
		var d byte
		switch hex[7] {
		case '0':
			d = 'R'
		case '1':
			d = 'D'
		case '2':
			d = 'L'
		case '3':
			d = 'U'
		default:
			panic("Can't map digit to direction")
		}

		v, _ := strconv.ParseInt(hex[2:7], 16, 0)

		return d, int(v)
	}
}

func solve18(lines []string, partTwo bool) {
	lastDir, _ := parse18(lines[len(lines)-1], partTwo)
	loc := Point{0, 0}

	acw, cw := 0, 0
	perimeter := 0
	grid := []Point{}

	for _, line := range lines {
		d, v := parse18(line, partTwo)
		perimeter += v

		if lastDir == 'R' && d == 'U' {
			acw++
		} else if lastDir == 'R' && d == 'D' {
			cw++
		} else if lastDir == 'L' && d == 'U' {
			cw++
		} else if lastDir == 'L' && d == 'D' {
			acw++
		} else if lastDir == 'U' && d == 'L' {
			acw++
		} else if lastDir == 'U' && d == 'R' {
			cw++
		} else if lastDir == 'D' && d == 'L' {
			cw++
		} else if lastDir == 'D' && d == 'R' {
			acw++
		} else {
			panic("Not a corner")
		}
		lastDir = d

		switch d {
		case 'R':
			loc = Point{loc.row, loc.col + v}
		case 'L':
			loc = Point{loc.row, loc.col - v}
		case 'U':
			loc = Point{loc.row - v, loc.col}
		case 'D':
			loc = Point{loc.row + v, loc.col}
		default:
			panic("Illegal direction")
		}
		grid = append(grid, loc)
	}

	if !(loc.row == 0 && loc.col == 0) {
		panic("Should be a closed loop")
	}

	corners := acw - cw
	if corners < 0 {
		corners = -corners
	}

	// Imagine the grid points correspond to the center of each sqaure.
	// Then the area inside the points can be found with the shoelace formula.
	// We must add in the area outside the points (to the edges of the squares).
	// For each straight line we expand out 0.5 to reach the edge (so add 0.5*perimeter).
	// Then add an extra 0.25 for each convex corner, and -0.25 for each concave one.
	fmt.Println(shoelace(grid) + (perimeter / 2) + (corners / 4))
}

// A = 0.5 * sum[i=1,N]{y[i] * (x[i-1] - x[i+1])}
func shoelace(points []Point) int {
	sum := 0
	for i, p := range points {
		yi := p.row
		ip := (i - 1)
		if ip < 0 {
			ip += len(points)
		}
		in := (i + 1)
		if in >= len(points) {
			in -= len(points)
		}
		xp := points[ip].col
		xn := points[in].col
		sum += yi * (xp - xn)
	}
	sum /= 2

	if sum > 0 {
		return sum
	}
	return -sum
}
