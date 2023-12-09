package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func day09() {
	file, err := os.Open("data/09.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var lines [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Split(line, " ")
		var ints []int
		for _, num := range nums {
			val, _ := strconv.Atoi(num)
			ints = append(ints, val)
		}
		lines = append(lines, ints)
	}

	partOne, partTwo := 0, 0
	for _, row := range lines {
		a, b := nextVals(row)
		partOne += a
		partTwo += b
	}

	fmt.Println(partOne)
	fmt.Println(partTwo)
}

func nextVals(row []int) (int, int) {
	rows := [][]int{row}

	// Build the rows of differences
	allZeroes := false
	for !allZeroes {
		allZeroes = true
		var nextRow []int
		prevRow := rows[len(rows)-1]
		for i := 1; i < len(prevRow); i++ {
			diff := prevRow[i] - prevRow[i-1]
			if diff != 0 {
				allZeroes = false
			}
			nextRow = append(nextRow, diff)
		}
		rows = append(rows, nextRow)
	}

	// Bubble up next val and return the answer
	start, end := 0, 0
	for i := len(rows) - 2; i >= 0; i-- {
		start = rows[i][0] - start
		end = end + rows[i][len(rows[i])-1]
	}

	return end, start
}
