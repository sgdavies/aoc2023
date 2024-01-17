package main

import (
	"fmt"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile("^(-?\\d+), (-?\\d+), (-?\\d+) @ +(-?\\d+), +(-?\\d+), +(-?\\d+)$")

type hailpath struct {
	px, py, pz int
	vx, vy, vz int
	m, c       float64 // y = m.x + c (ignoring z)
}

func parseHailstone(line string) hailpath {
	matches := re.FindSubmatch([]byte(line))
	px, _ := strconv.Atoi(string(matches[1]))
	py, _ := strconv.Atoi(string(matches[2]))
	pz, _ := strconv.Atoi(string(matches[3]))
	vx, _ := strconv.Atoi(string(matches[4]))
	vy, _ := strconv.Atoi(string(matches[5]))
	vz, _ := strconv.Atoi(string(matches[6]))

	// y = m.x + c
	// m = vy/vx
	// x = px + t.vx and y = py + t.vy
	// y=c when x=0; x=0 when t= -px/vx
	// => c = py+t.vy == py-(px/vx).vy == py-px.(vy/vx) == py -px.m
	m := float64(vy) / float64(vx)
	c := float64(py) - float64(px)*m

	return hailpath{px, py, pz, vx, vy, vz, m, c}
}

func day24() {
	test := false
	var file string
	var start, end int
	if test {
		file = "data/24.ex"
		start, end = 7, 27
	} else {
		file = "data/24.txt"
		start, end = 200000000000000, 400000000000000
	}
	lines := LinesFromFile(file)
	var paths []hailpath
	for _, line := range lines {
		paths = append(paths, parseHailstone(line))
	}

	ans := 0
	for i, a := range paths[:len(paths)-1] {
		for _, b := range paths[i+1:] {
			if intersect(a, b, start, end) {
				ans += 1
			}
		}
	}

	fmt.Println(ans)
}

func intersect(a, b hailpath, start, end int) bool {
	lo, hi := float64(start), float64(end)
	if a.m == b.m {
		if a.c == b.c {
			panic("Haven't implemented identical paths")
		} else {
			return false // Parallel
		}
	}

	x := (b.c - a.c) / (a.m - b.m)
	y := a.m*x + a.c

	// Check - but properly needs math.Abs(y - y2) < epsilon
	// y2 := b.m*x + b.c
	// if y != y2 {
	// 	panic(fmt.Sprintf("Should be equal: %v %v", y, y2))
	// }

	if x < lo || x > hi || y < lo || y > hi {
		return false
	}

	ta := (x - float64(a.px)) / float64(a.vx)
	tb := (x - float64(b.px)) / float64(b.vx)

	if ta < 0 || tb < 0 {
		return false
	}

	return true
}
