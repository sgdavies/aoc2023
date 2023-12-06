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

func day04() {
	file, err := os.Open("data/04.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var matches []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches = append(matches, scratch(line))
	}

	partOne := 0
	for _, s := range matches {
		partOne += int(math.Pow(2, float64(s-1)))
	}
	fmt.Println(partOne)

	var allCards []int
	for range matches {
		allCards = append(allCards, 1) // Original card
	}
	partTwo := 0 // Total number of scratch cards
	for i := range matches {
		partTwo += allCards[i]

		if matches[i] > 0 {
			for j := 1; j <= matches[i]; j++ {
				allCards[i+j] += allCards[i]
			}
		}
	}
	fmt.Println(partTwo)
}

func scratch(line string) int {
	// "Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1"
	// "Card 3:", "  1 21 53 59 44 | 69 82 63 72 16 21 14  1"
	// "1 21 53 59 44 | 69 82 63 72 16 21 14  1"
	// "1 21 53 59 44", "69 82 63 72 16 21 14  1"
	// "1 21 53 59 44", "69 82 63 72 16 21 14 1"
	// 1, 21, 53, 59, 44 ; ...
	// 1, 21, 44, 53, 59
	parts := strings.Split(line, ":")
	part2 := strings.TrimSpace(parts[1])
	partNumbers := strings.Split(part2, " | ")
	winningNumbers := scratchNumbers(partNumbers[0])
	myNumbers := scratchNumbers(partNumbers[1])

	score := 0
	for _, mine := range myNumbers {
		if slices.Contains(winningNumbers, mine) {
			score += 1
		}
	}

	return score
}

func scratchNumbers(line string) []int {
	// Remove double spaces so Split(" ") doesn't create empty entries
	line = strings.ReplaceAll(line, "  ", " ")
	numStrs := strings.Split(line, " ")

	var nums []int
	for _, s := range numStrs {
		val, _ := strconv.Atoi(s)
		nums = append(nums, val)
	}

	slices.Sort(nums)
	return nums
}
