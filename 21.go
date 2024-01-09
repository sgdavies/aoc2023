package main

import (
	"fmt"
	"math"
	"slices"
)

func day21() {
	lines := LinesFromFile("data/21.txt")
	garden := make(map[Point]bool)
	var start Point
	for row, line := range lines {
		for col, r := range line {
			p := Point{row, col}
			if r == 'S' {
				start = p
				// locs[p] = true
				garden[p] = true
			} else if r == '.' {
				garden[p] = true
			}
		}
	}
	// fmt.Println(start)
	// fmt.Println(garden[Point{128, 128}])

	partOneSteps := 64
	ans := iterate21(start, garden, []int{partOneSteps, 129, 130, 160, 161}, false)
	// fmt.Println(ans)
	fmt.Println(ans[partOneSteps])

	// solve21PartTwo(start, garden, ans, 3)  //65)
	// solve21PartTwo(start, garden, ans, 10) //196)
	// solve21PartTwo(start, garden, ans, 17) //327)
	// solve21PartTwo(start, garden, ans, 24) //589)
	// solve21PartTwo(start, garden, ans, 31) //720)
	solve21PartTwoQuadFormula(start, garden)
}

func solve21PartTwoQuadFormula(start Point, garden map[Point]bool) {
	// (1) Use updated solve function to actually fill out the infinite grid
	// (2) use that to extrapolate to answer from first 3 edge points
	//     (65, 65+131, 65+131*131)

	// ans = a*steps^2 + b*steps + c
	// (1) a*x1^2 + b*x1 + c = R1
	// (2) a*x2^2 + b*x2 + c = R2
	// (3) a*x3^2 + b*x3 + c = R3
	// (1)*x2^2 - (2)*x1^2 ==> (4) b(x1*x2^2 - x2*x1^2) + c(x2^2 - x1^2) = R1*x2^2 - R2*x1^2
	// (1)*x3^2 - (3)*x1^2 ==> (5) b(x1*x3^2 - x3*x1^2) + c(x3^2 - x1^2) = R1*x3^2 - R3*x1^2
	// z12 = x1*x2^2 - x2*x1^2, z13 = x1*x3^2 - x3*x1^2
	// (4)*z13 - (5)*z12 ==> (6) c(z13(x2^2-x1^2) - z12(x3^2 - x1^2) = z13(R1*x2^2 - R2*x1^2) - z12(R1*x3^2 - R3*x1^2)
	// => c [from (6)]
	// => b [from (4)]
	// => a [from (1)]
	xi1 := 65
	xi2 := 65 + 131
	xi3 := 65 + 2*131
	results := iterate21(start, garden, []int{64, xi1, xi2, xi3}, true)
	x1, x2, x3 := float64(xi1), float64(xi2), float64(xi3)
	R1 := float64(results[xi1])
	R2 := float64(results[xi2])
	R3 := float64(results[xi3])

	z12 := x1*x2*x2 - x2*x1*x1
	z13 := x1*x3*x3 - x3*x1*x1
	c := (z13*(R1*x2*x2-R2*x1*x1) - z12*(R1*x3*x3-R3*x1*x1)) / (z13*(x2*x2-x1*x1) - z12*(x3*x3-x1*x1))
	b := ((R1*x2*x2 - R2*x1*x1) - c*(x2*x2-x1*x1)) / (x1*x2*x2 - x2*x1*x1)
	a := (R1 - c - b*x1) / (x1 * x1)

	// fmt.Println(a, b, c)
	// fmt.Println(x1, x2, x3)
	// fmt.Println(R1, R2, R3)
	// fmt.Println(a*x1*x1+b*x1+c, a*x2*x2+b*x2+c, a*x3*x3+b*x3+c)
	x := float64(26501365)
	fmt.Printf("%.0f\t%.0f\n", x, a*x*x+b*x+c)
	// x = float64(589)
	// fmt.Printf("%v\t%.0f\n", x, a*x*x+b*x+c)
	// x = float64(720)
	// fmt.Printf("%v\t%.0f\n", x, a*x*x+b*x+c)
}

func solve21PartTwo(start Point, garden map[Point]bool, ans map[int]int, target int) {
	// Part two
	fmt.Println("Target: ", target)
	n := 2*start.col + 1 // Width of each tile
	s := target          //26501365
	k := 10              // 161             //130    // Number of steps by which we're in steady state (even=7457, odd=7520)
	sk := s - k          // the tile that ends before this is guaranteed to be completely full

	// Calculate number of tiles (tN) and number of steps (sN) for the edge of the first
	// completed tile (t0) and the next 2 tiles (2 because 2*n > k)
	t0 := int(math.Floor((((2*float64(sk) - 1) / float64(n)) + 1) / 2))
	s0 := tToS(t0, n)
	t1 := t0 + 1
	s1 := tToS(t1, n)
	t2 := t1 + 1
	s2 := tToS(t2, n)

	// Now calculate the number of steps to reach the corner of tiles not on the axes
	sc0 := tToSc(t0, n)
	sc1 := tToSc(t1, n)
	sc2 := tToSc(t2, n)

	// fmt.Println(t0, s0, sc0, t1, s1, sc1, t2, s2, sc2)

	// The answer is the answer for each of the completely filled tiles (odd or even),
	// plus the part-filled tiles on the axes (filled from a side point)
	// plus the part-fille tiles on the diagonal edges (filled from the corner)
	s0d := s - s0
	s1d := s - s1 // if -ve, ignore
	s2d := s - s2
	// fmt.Println(s0d, s1d, s2d)
	axesLeft := iterate21(Point{start.row, 0}, garden, []int{s0d, s1d, s2d}, false)
	axesRight := iterate21(Point{start.row, n}, garden, []int{s0d, s1d, s2d}, false)
	axesTop := iterate21(Point{0, start.col}, garden, []int{s0d, s1d, s2d}, false)
	axesBottom := iterate21(Point{n, start.col}, garden, []int{s0d, s1d, s2d}, false)

	// part-filled diagonals
	sc0d := s - sc0
	sc1d := s - sc1
	sc2d := s - sc2
	fmt.Println(sc0d, sc1d, sc2d)
	diagTopLeft := iterate21(Point{0, 0}, garden, []int{sc0d, sc1d, sc2d}, false)
	diagTopRight := iterate21(Point{0, n}, garden, []int{sc0d, sc1d, sc2d}, false)
	diagBottomLeft := iterate21(Point{n, 0}, garden, []int{sc0d, sc1d, sc2d}, false)
	diagBottomRight := iterate21(Point{n, n}, garden, []int{sc0d, sc1d, sc2d}, false)

	tFilled := t0 - 1
	evenFilledCount := evenFilled(tFilled)
	oddFilledCount := oddFilled(tFilled)
	evenFilledCount *= evenFilledCount
	oddFilledCount *= oddFilledCount
	fmt.Println(tFilled, t0, t1, t2, evenFilledCount, oddFilledCount)

	// Calculate the answer!
	// The completely-filled tiles (even and odd)
	cfe := evenFilledCount * ans[160]
	cfo := oddFilledCount * ans[161]
	fmt.Println("\t", cfe, cfo)

	// One each of the unfilled squares on each axis in tile t0, t1, t2
	cal := axesLeft[s0d] + axesLeft[s1d] + axesLeft[s2d] // If not present, will get 0
	car := axesRight[s0d] + axesRight[s1d] + axesRight[s2d]
	cat := axesTop[s0d] + axesTop[s1d] + axesTop[s2d]
	cab := axesBottom[s0d] + axesBottom[s1d] + axesBottom[s2d]
	fmt.Println("\t", cal, car, cat, cab)

	// Unfilled diagonals. Each diagonal at scN is tN tiles long
	cd0 := t0 * (diagTopLeft[sc0d] + diagTopRight[sc0d] + diagBottomLeft[sc0d] + diagBottomRight[sc0d])
	cd1 := t1 * (diagTopLeft[sc1d] + diagTopRight[sc1d] + diagBottomLeft[sc1d] + diagBottomRight[sc1d])
	cd2 := t2 * (diagTopLeft[sc2d] + diagTopRight[sc2d] + diagBottomLeft[sc2d] + diagBottomRight[sc2d])
	fmt.Println("\t", cd0, cd1, cd2)

	fmt.Println(cfe + cfo + cal + car + cat + cab + cd0 + cd1 + cd2)
}

func tToS(t, n int) int {
	if t < 1 {
		return 0
	}
	return ((2*t-1)*n + 1) / 2
}

func tToSc(t, n int) int {
	if t < 0 {
		return 0
	}
	return t*n + 1
}

func evenFilled(tFilled int) int {
	if tFilled < 0 {
		return 0
	}
	return (2 * int(math.Floor(float64(tFilled)/2))) + 1
}

func oddFilled(tFilled int) int {
	if tFilled < 0 {
		return 0
	}
	return 2 * int(math.Floor((float64(tFilled)+1)/2))
}

// Calculate number of reachable squares, at each of the passed-in step counts
func iterate21(start Point, garden map[Point]bool, steps []int, partTwo bool) map[int]int {
	locs := make(map[Point]bool)
	locs[start] = true
	n := 2*start.col + 1 // used for part 2 calcs

	max := steps[0]
	for _, m := range steps[1:] {
		if m > max {
			max = m
		}
	}

	ret := make(map[int]int)

	for s := 1; s <= max; s++ {
		newLocs := make(map[Point]bool)
		for loc := range locs {
			for _, delta := range []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				np := Point{loc.row + delta.row, loc.col + delta.col}
				if !partTwo {
					if garden[np] {
						newLocs[np] = true
					}
				} else {
					normNp := Point{norm21(np.row, n), norm21(np.col, n)}
					if garden[normNp] {
						newLocs[np] = true
					}
				}
			}
		}
		locs = newLocs
		if slices.Contains(steps, s) {
			ret[s] = len(locs)
		}
	}

	return ret
}

func norm21(c int, size int) int {
	// Normalize grid coord to 0 <= g <= size
	if c < 0 {
		tiles := -(c / size)
		c += (tiles + 1) * size
		if c < 0 {
			panic("add more!")
		}
	}
	return c % size
}
