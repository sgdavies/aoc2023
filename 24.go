package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

var re = regexp.MustCompile("^(-?\\d+), (-?\\d+), (-?\\d+) @ +(-?\\d+), +(-?\\d+), +(-?\\d+)$")

type hailpath struct {
	p    []int   // px, py, pz
	v    []int   // vx, vy, vz
	m, c float64 // y = m.x + c (ignoring z)
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

	return hailpath{[]int{px, py, pz}, []int{vx, vy, vz}, m, c}
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

	// Part two...
	ps, vs := [][]float64{}, [][]float64{}
	for h := 0; h < 5; h++ {
		pp, vv := []float64{}, []float64{}
		for axis := 0; axis < 3; axis++ {
			pp = append(pp, float64(paths[h].p[axis]))
			vv = append(vv, float64(paths[h].v[axis]))
		}
		ps = append(ps, pp)
		vs = append(vs, vv)
	}
	// This should work, but doesn't:
	// Pnx + Vnx.tn = Px + Vx.tn => tn = (Pnx-Px)/(Vx-Vnx) = (Pny-Py)/(Vy-Vny) = (Pnz-Pz)/(Vz-Vnz)
	// Multiply out and sub in values for H1 and H2 to get 4 simul eqs with 4 unknowns
	// Solve for these unknowns and then use these to calc the remaining 2 unknowns.
	// a := mat.NewDense(4, 4, []float64{
	// 	vs[0][1], -vs[0][0], -ps[0][1], ps[0][0],
	// 	vs[1][1], -vs[1][0], -ps[1][1], ps[1][0],
	// 	vs[2][1], -vs[2][0], -ps[2][1], ps[2][0],
	// 	vs[3][1], -vs[3][0], -ps[3][1], ps[3][0]})
	// b := mat.NewDense(4, 1, []float64{
	// 	ps[0][0]*vs[0][1] - ps[0][1]*vs[0][0],
	// 	ps[1][0]*vs[1][1] - ps[1][1]*vs[1][0],
	// 	ps[2][0]*vs[2][1] - ps[2][1]*vs[2][0],
	// 	ps[3][0]*vs[3][1] - ps[3][1]*vs[3][0]})
	// fmt.Println(a)
	// fmt.Println(b)
	// var x mat.Dense
	// err := x.Solve(a, b)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(x)
	// }

	// Pn + Vn.tn = P + V.tn => (Pn-P) = -tn.(Vn-V)
	// tn is scalar, so vectors (Pn-P) and (Vn-V) are parallel
	// Therefore cross-product is zero.
	// When multiplied out, the non-linear terms (PaVb) cancel out.
	// Then compare & subtract from different hailstones to get 6 simultaneous equations
	// for the 6 unknown variables:
	// v0y-v1y  v1x-v0x  0        p1y-p0y  p0x-p1x  0       = p0x*v0y - p0y*v0x - p1x*v1y + p1y*v1x
	// v0y-v2y  v2x-v0x  0        p2y-p0y  p0x-p2x  0       = p0x*v0y - p0y*v0x - p2x*v2y + p2y*v2x
	// v0y-v3y  v3x-v0x  0        p3y-p0y  p0x-p3x  0       = p0x*v0y - p0y*v0x - p3x*v3y + p3y*v3x
	// v0z-v1z  0        v1x-v0x  p1z-p0z  0        p0x-p1x = p0x*v0z - p0z*v0x - p1x*v1z + p1z*v1x
	// v0z-v2z  0        v2x-v0x  p2z-p0z  0        p0x-p2x = p0x*v0z - p0z*v0x - p2x*v2z + p2z*v2x
	// v0z-v3z  0        v3x-v0x  p3z-p0z  0        p0x-p3x = p0x*v0z - p0z*v0x - p3x*v3z + p3z*v3x
	// Solve as A.x = B
	a := mat.NewDense(6, 6, []float64{
		vs[0][1] - vs[1][1], vs[1][0] - vs[0][0], 0, ps[1][1] - ps[0][1], ps[0][0] - ps[1][0], 0,
		vs[0][1] - vs[2][1], vs[2][0] - vs[0][0], 0, ps[2][1] - ps[0][1], ps[0][0] - ps[2][0], 0,
		vs[0][1] - vs[3][1], vs[3][0] - vs[0][0], 0, ps[3][1] - ps[0][1], ps[0][0] - ps[3][0], 0,
		vs[0][2] - vs[1][2], 0, vs[1][0] - vs[0][0], ps[1][2] - ps[0][2], 0, ps[0][0] - ps[1][0],
		vs[0][2] - vs[2][2], 0, vs[2][0] - vs[0][0], ps[2][2] - ps[0][2], 0, ps[0][0] - ps[2][0],
		vs[0][2] - vs[3][2], 0, vs[3][0] - vs[0][0], ps[3][2] - ps[0][2], 0, ps[0][0] - ps[3][0],
	})
	b := mat.NewDense(6, 1, []float64{
		ps[0][0]*vs[0][1] - ps[0][1]*vs[0][0] - ps[1][0]*vs[1][1] + ps[1][1]*vs[1][0],
		ps[0][0]*vs[0][1] - ps[0][1]*vs[0][0] - ps[2][0]*vs[2][1] + ps[2][1]*vs[2][0],
		ps[0][0]*vs[0][1] - ps[0][1]*vs[0][0] - ps[3][0]*vs[3][1] + ps[3][1]*vs[3][0],
		ps[0][0]*vs[0][2] - ps[0][2]*vs[0][0] - ps[1][0]*vs[1][2] + ps[1][2]*vs[1][0],
		ps[0][0]*vs[0][2] - ps[0][2]*vs[0][0] - ps[2][0]*vs[2][2] + ps[2][2]*vs[2][0],
		ps[0][0]*vs[0][2] - ps[0][2]*vs[0][0] - ps[3][0]*vs[3][2] + ps[3][2]*vs[3][0],
	})
	var xxx mat.Dense
	err := xxx.Solve(a, b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(xxx)
	// pxi := xxx.At(0, 0)
	// pyi := xxx.At(1, 0)
	// pzi := xxx.At(2, 0)
	// fmt.Println(px + py + pz) // doesn't work - presume nums are too big to round accurately

	// Instead, get V (which are smaller so are definitely correct) and extrapolate back to find P
	vx := int(math.Round(xxx.At(3, 0)))
	vy := int(math.Round(xxx.At(4, 0)))
	vz := int(math.Round(xxx.At(5, 0)))

	// M0 = (v0y-vy)/(v0x-vx)
	// M1 = (v1y-vy)/(v1x-vx)
	m0 := (vs[0][1] - float64(vy)) / (vs[0][0] - float64(vx))
	m1 := (vs[1][1] - float64(vy)) / (vs[1][0] - float64(vx))
	// C0 = p0y - (M0*p0x)
	// C1 = p1y - (M1*p1x)
	c0 := ps[0][1] - m0*ps[0][0]
	c1 := ps[1][1] - m1*ps[1][0]
	// Px = int((C1-C0)/(M0-M1))
	// Py = int(MA*Px + C0)
	// time = (Px - p0x)//(v0x-vx)
	// Pz = p0z + (v0z - vz)*time
	px := (c1 - c0) / (m0 - m1)
	py := m0*px + c0
	t := (px - ps[0][0]) / (vs[0][0] - float64(vx))
	pz := ps[0][2] + (vs[0][2]-float64(vz))*t

	total := px + py + pz
	fmt.Println( /*total, math.Round(total),*/ int(math.Round(total)))
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

	if x < lo || x > hi || y < lo || y > hi {
		return false
	}

	ta := (x - float64(a.p[0])) / float64(a.v[0])
	tb := (x - float64(b.p[0])) / float64(b.v[0])

	if ta < 0 || tb < 0 {
		return false
	}

	return true
}
