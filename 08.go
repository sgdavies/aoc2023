package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type node struct {
	left  string
	right string
}

func day08() {
	file, err := os.Open("data/08.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	turns := lines[0]
	nodes := make(map[string]node)
	for _, line := range lines[2:] {
		parts := strings.Split(line, " = ")
		lrParts := strings.Split(parts[1], ", ")
		nodes[parts[0]] = node{lrParts[0][1:], lrParts[1][:len(lrParts[1])-1]} // Strip off "(" and ")"
	}

	// Part one
	steps := 0
	nodeName := "AAA"
	for nodeName != "ZZZ" {
		if turns[steps%len(turns)] == 'L' {
			nodeName = nodes[nodeName].left
		} else {
			nodeName = nodes[nodeName].right
		}
		steps += 1
	}

	fmt.Println(steps)

	// Part two
	// The naïve solution would be to run all at once
	// until the end condition - but this takes far too long.
	// Instead find the individual lengths and then the LCM.
	var ghostSteps []int
	for k, _ := range nodes {
		if k[len(k)-1] == 'A' {
			// Run until we hit a ..Z
			steps = 0
			nodeName = k

			for nodeName[len(nodeName)-1] != 'Z' {
				if turns[steps%len(turns)] == 'L' {
					nodeName = nodes[nodeName].left
				} else {
					nodeName = nodes[nodeName].right
				}
				steps += 1
			}
			ghostSteps = append(ghostSteps, steps)
		}
	}

	// fmt.Println(ghostSteps)
	// fmt.Println(GCD(48, 18))
	// fmt.Println(LCM([]int{8, 9, 21}))
	fmt.Println(LCM(ghostSteps))
}

func GCD(a, b int) int {
	if a < b {
		a, b = b, a
	}

	if b == 0 {
		return a
	} else {
		return GCD(b, a%b)
	}
}

func LCM(ints []int) int {
	if len(ints) < 2 {
		panic("Need at least 2 ints!")
	}
	a, b := ints[0], ints[1]
	lcmAB := a * (b / GCD(a, b))

	newInts := append([]int{lcmAB}, ints[2:]...)
	if len(newInts) == 1 {
		return lcmAB
	} else {
		return LCM(newInts)
	}
}
