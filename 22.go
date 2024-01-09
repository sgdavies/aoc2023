package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type P3 struct {
	x int
	y int
	z int
}

type block struct {
	points     []P3
	restingOn  map[*block]bool
	supporting map[*block]bool
}

func day22() {
	world := make(map[P3]*block) // Fills up as the blocks come to rest
	lines := LinesFromFile("data/22.txt")
	blocks := []*block{}
	for _, line := range lines {
		block := parseBlock(line)
		blocks = append(blocks, &block)
	}
	slices.SortFunc(blocks, func(a, b *block) int { return a.lowest() - b.lowest() })

	for _, b := range blocks {
		// TODO: ensure blocks keep falling (need another loop?)
		for {
			if b.lowest() == 1 {
				// z=1 means we are resting on the floor (at z=0) so stop falling
				for _, p := range b.points {
					world[p] = b
				}
				break // Next block
			}

			nextPoints := []P3{}
			for _, p := range b.points {
				nextPoints = append(nextPoints, P3{p.x, p.y, p.z - 1})
			}
			for _, np := range nextPoints {
				if otherBlock, pres := world[np]; pres {
					b.restingOn[otherBlock] = true
					otherBlock.supporting[b] = true
				}
			}
			if len(b.restingOn) > 0 {
				for _, p := range b.points {
					world[p] = b
				}
				break // Next block
			} else {
				b.points = nextPoints
			}
		}
	}

	ans := 0
	for _, b := range blocks {
		// For part one, count any blocks that could be removed
		// That means they're not supporting any other blocks, or
		// for each block they are supporting, that block is resting
		// on at least one other block.
		if len(b.supporting) == 0 {
			ans++
		} else {
			allMultiplySupported := true
			for ob, _ := range b.supporting {
				if len(ob.restingOn) < 2 {
					allMultiplySupported = false
					break
				}
			}

			if allMultiplySupported {
				ans++
			}
		}
	}

	fmt.Println(ans)

	// Part two
	partTwo := 0
	for _, b := range blocks {
		otherCount := 0
		othersToFall := []*block{b}
		fallen := make(map[*block]bool)
		for len(othersToFall) > 0 {
			nb := othersToFall[0]
			othersToFall = othersToFall[1:]
			fallen[nb] = true
			for ob := range nb.supporting {
				stillSupported := false
				for rob := range ob.restingOn {
					if !fallen[rob] {
						stillSupported = true
						break
					}
				}
				if !stillSupported {
					otherCount++
					othersToFall = append(othersToFall, ob)
				}
			}
		}
		partTwo += otherCount
	}
	fmt.Println(partTwo)
}

func (b *block) lowest() int {
	lowZ := b.points[0].z
	for _, p := range b.points[1:] {
		if p.z < lowZ {
			lowZ = p.z
		}
	}
	return lowZ
}

func parseBlock(line string) block {
	parts := strings.Split(line, "~")
	start := strings.Split(parts[0], ",")
	end := strings.Split(parts[1], ",")
	si := nums(start)
	ei := nums(end)
	sp := P3{si[0], si[1], si[2]}
	ep := P3{ei[0], ei[1], ei[2]}

	var d P3
	var length int
	if sp == ep {
		d = P3{0, 0, 0}
		length = 0
	} else if sp.x < ep.x {
		d = P3{1, 0, 0}
		length = ep.x - sp.x
	} else if sp.y < ep.y {
		d = P3{0, 1, 0}
		length = ep.y - sp.y
	} else if sp.z < ep.z {
		d = P3{0, 0, 1}
		length = ep.z - sp.z
	} else {
		panic("Implement backwards blocks")
	}

	b := block{[]P3{sp}, make(map[*block]bool), make(map[*block]bool)}
	dp := sp
	for i := 0; i < length; i++ {
		dp = P3{dp.x + d.x, dp.y + d.y, dp.z + d.z}
		b.points = append(b.points, dp)
	}

	return b
}

func nums(in []string) []int {
	out := []int{}
	for _, s := range in {
		v, _ := strconv.Atoi(s)
		out = append(out, v)
	}
	return out
}
