package main

import (
	"bufio"
	"fmt"
	"os"
)

func day11() {
	file, err := os.Open("data/11.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	populatedRows := []bool{}
	populatedCols := []bool{}
	galaxies := []Point{}

	scanner := bufio.NewScanner(file)
	row := 0
	for scanner.Scan() {
		populatedRows = append(populatedRows, false)
		for c, r := range scanner.Text() {
			if c >= len(populatedCols) { // Only on first time through the loop
				populatedCols = append(populatedCols, false)
			}

			if r == '#' {
				galaxies = append(galaxies, Point{row: row, col: c})
				populatedRows[row] = true
				populatedCols[c] = true
			}
		}
		row += 1
	}

	fmt.Println(solveGalaxies(galaxies, populatedCols, populatedRows, 1))
	fmt.Println(solveGalaxies(galaxies, populatedCols, populatedRows, 1000000-1)) // Don't count the original row/column twice
}

func solveGalaxies(galaxies []Point, popCols []bool, popRows []bool, age int) int {
	answer := 0
	for i, ga := range galaxies[:len(galaxies)-1] {
		for _, gb := range galaxies[i+1:] {
			answer += galaxyDist(ga, gb, popCols, popRows, age)
		}
	}

	return answer
}

func galaxyDist(a Point, b Point, popCols []bool, popRows []bool, age int) int {
	var h_start, h_end, v_start, v_end int

	if a.col < b.col {
		h_start, h_end = a.col, b.col
	} else {
		h_start, h_end = b.col, a.col
	}

	if a.row < b.row {
		v_start, v_end = a.row, b.row
	} else {
		v_start, v_end = b.row, a.row
	}

	h_gaps := 0
	for _, b := range popCols[h_start:h_end] {
		if !b {
			h_gaps += 1
		}
	}
	v_gaps := 0
	for _, b := range popRows[v_start:v_end] {
		if !b {
			v_gaps += 1
		}
	}

	return (h_end - h_start + v_end - v_start) + age*(h_gaps+v_gaps)
}
