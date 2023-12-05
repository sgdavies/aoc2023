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

type rrange struct {
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
	maps := [][]rrange{se2so, so2fe, fe2wa, wa2li, li2te, te2hu, hu2lo}

	// Part one
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

	// Part two
}

// Returns mappings, sorted by source
func readPlantMap(scanner *bufio.Scanner, name string) []rrange {
	scanner.Scan()
	line := scanner.Text() // x-to-y map:
	if !(strings.HasPrefix(line, name)) {
		panic(line + " should be " + name)
	}

	var plantMap []rrange
	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			break
		}
		nums := strings.Split(line, " ")
		dest, _ := strconv.Atoi(nums[0])
		source, _ := strconv.Atoi(nums[1])
		length, _ := strconv.Atoi(nums[2])
		plantMap = append(plantMap, rrange{source, dest, length})
	}

	slices.SortFunc(plantMap, func(a rrange, b rrange) int { return a.source - b.source })
	return plantMap
}

func plantLookup(k int, rr []rrange) int {
	// TODO
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
