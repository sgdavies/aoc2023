package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type rangeMap struct {
	source int // Sort based on this
	dest   int
	length int
}

func day05() {
	file, err := os.Open("data/05.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	var seeds []int
	for _, seed := range strings.Split(line[7:], " ") {
		v, _ := strconv.Atoi(seed)
		seeds = append(seeds, v)
	}

	scanner.Scan()
	line = scanner.Text() // blank
	se2so := readPlantMap(scanner, "seed-to-soil map")
	so2fe := readPlantMap(scanner, "soil-to-fertilizer map")
	fe2wa := readPlantMap(scanner, "fertilizer-to-water map")
	wa2li := readPlantMap(scanner, "water-to-light map")
	li2te := readPlantMap(scanner, "light-to-temperature map")
	te2hu := readPlantMap(scanner, "temperature-to-humidity map")
	hu2lo := readPlantMap(scanner, "humidity-to-location map")
	maps := [][]rangeMap{se2so, so2fe, fe2wa, wa2li, li2te, te2hu, hu2lo}

	// Part one
	day05partOne(seeds, maps)

	// Part two
	day05partTwo(seeds, maps)
	day05partTwoBad(seeds, maps)
}

func day05partOne(seeds []int, maps [][]rangeMap) {
	bestLocation := math.MaxInt
	for _, seed := range seeds {
		next := seed
		for _, lookup := range maps {
			next = plantLookup(next, lookup)
		}
		if next < bestLocation {
			bestLocation = next
		}
	}
	fmt.Println(bestLocation)
}

func day05partTwoBad(seeds []int, maps [][]rangeMap) {
	// brute force :(
	best := math.MaxInt
	for i := 0; i < len(seeds); i += 2 {
		fmt.Printf("\t%d of %d\n", i, len(seeds))
		start := seeds[i]
		length := seeds[i+1]
		for j := start; j < start+length; j++ {
			next := j
			for _, lookup := range maps {
				next = plantLookup(next, lookup)
			}
			if next < best {
				best = next
			}
		}
	}

	fmt.Println(best)
}

// There's a bug in here somewhere that only hits the real data, not the test data
// Given the brute force version took 2 minutes to code and 4 to run, I'm stopping debugging here.
func day05partTwo(seeds []int, maps [][]rangeMap) {
	var numberRanges []intRange
	for i := 0; i < len(seeds); i += 2 {
		numberRanges = append(numberRanges, intRange{seeds[i], seeds[i] + seeds[i+1] - 1}) // -1 as we count the first number as well
	}
	// TODO - for testing, replace with {a,a}, {b,b}, {c,c}, ... instead of {a,a+b}, {c,c+d}, ... - should get part one answer
	// numberRanges = []intRange{{79, 92}, {55, 67}}
	for _, lookup := range maps {
		var nextNumberRanges []intRange
		lx := 0 // lookup index
		lu := lookup[lx]

		for len(numberRanges) > 0 { // Consume all the input
			r := numberRanges[0]
			if r.end < lu.source {
				// Entire range misses - direct mapping
				nextNumberRanges = append(nextNumberRanges, r)
				numberRanges = numberRanges[1:] // input range fully consumed
			} else if r.start < lu.source {
				// Some overlap, so we'll need to split input. Determine which type of overlap.
				// s1---e1     or  s1----e1              or  s1-------e1
				//    s2---e2        s2--e2 (exact end)         s2-e2
				if r.end <= lu.source+lu.length {
					nextNumberRanges = append(nextNumberRanges, intRange{r.start, lu.source - 1})
					nextNumberRanges = append(nextNumberRanges, intRange{lu.dest, r.end + lu.dest - lu.source})
					// input range fully consumed
					numberRanges = numberRanges[1:]

					if r.end == lu.source+lu.length {
						// lookup range fully consumed
						lx, lu = updateLookup(lx, lookup)
					}
				} else {
					// input range wider than lookup range
					nextNumberRanges = append(nextNumberRanges, intRange{r.start, lu.source - 1})
					nextNumberRanges = append(nextNumberRanges, intRange{lu.dest, lu.dest + lu.length})

					// Add remaining bit of input back into numberRanges
					numberRanges = append([]intRange{{lu.source + lu.length + 1, r.end}}, numberRanges[1:]...)

					// lookup range fully consumed
					// TODO: worry about overlapping ranges
					lx, lu = updateLookup(lx, lookup)
				}
			} else if r.start <= lu.source+lu.length {
				// Some overlap, so we'll need to split input. Determine which type of overlap.
				//   s1-e1    or     s1--e1              or     s1---e1
				// s2-----e2      s2-----e2 (exact end)      s2---e2
				if r.end <= lu.source+lu.length {
					nextNumberRanges = append(nextNumberRanges, intRange{r.start + lu.dest - lu.source, r.end + lu.dest - lu.source})
					if r.end == lu.source+lu.length {
						lx, lu = updateLookup(lx, lookup) // lookup range fully consumed
					}
					numberRanges = numberRanges[1:] // input range fully consumed
				} else {
					// Split the input
					nextNumberRanges = append(nextNumberRanges, intRange{r.start + lu.dest - lu.source, lu.dest + lu.length})
					// Add remaining bit of input back into numberRanges
					numberRanges = append([]intRange{{lu.source + lu.length + 1, r.end}}, numberRanges[1:]...)
					lx, lu = updateLookup(lx, lookup) // lookup range fully consumed
				}
			} else {
				// We've gone past the previous lookup. Update it (the sentinal is INTMAX so we can never go past that)
				lx, lu = updateLookup(lx, lookup)
			}
		}

		slices.SortFunc(nextNumberRanges, func(a, b intRange) int { return a.start - b.start })
		numberRanges = mergeRanges(nextNumberRanges)
	}

	fmt.Println(numberRanges[0].start)
}

// Returns mappings, sorted by source
func readPlantMap(scanner *bufio.Scanner, name string) []rangeMap {
	scanner.Scan()
	line := scanner.Text() // x-to-y map:
	if !(strings.HasPrefix(line, name)) {
		panic(line + " should be " + name)
	}

	var plantMap []rangeMap
	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			break
		}
		nums := strings.Split(line, " ")
		dest, _ := strconv.Atoi(nums[0])
		source, _ := strconv.Atoi(nums[1])
		length, _ := strconv.Atoi(nums[2])
		plantMap = append(plantMap, rangeMap{source, dest, length - 1}) // Length-1 as we also count the first entry
	}

	slices.SortFunc(plantMap, func(a rangeMap, b rangeMap) int { return a.source - b.source })
	return plantMap
}

func plantLookup(k int, rr []rangeMap) int {
	for _, r := range rr {
		if k < r.source {
			// No match
			return k
		} else if k <= r.source+r.length {
			return k - r.source + r.dest
		}
	}
	return k
}

// Part two
type intRange struct {
	start int
	end   int
}

func updateLookup(lx int, lookup []rangeMap) (int, rangeMap) {
	// TODO: worry about overlapping ranges
	lx += 1
	if lx < len(lookup) {
		return lx, lookup[lx]
	} else {
		return lx, rangeMap{math.MaxInt64, math.MaxInt64, 1} // hack
	}
}

// Assumes range is sorted on elem.start
func mergeRanges(in []intRange) []intRange {
	changed := true
	for changed {
		changed = false
		for i := range in {
			if i == len(in)-1 {
				break
			}

			if in[i].end+1 >= in[i+1].start { // +1 - contiguous ints
				before := append(in[:i], intRange{in[i].start, max(in[i].end, in[i+1].end)})
				in = append(before, in[i+2:]...)
				changed = true
				break
			}
		}
	}

	return in
}
