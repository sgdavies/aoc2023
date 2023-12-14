package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

var (
	dishHeight int
	dishWidth  int
)

func day14() {
	file, err := os.Open("data/14.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	dish := make(map[Point]rune)
	row := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for col, r := range line {
			if col > dishWidth {
				dishWidth = col
			}

			if r == 'O' || r == '#' {
				dish[Point{row, col}] = r
			}
		}

		row += 1
	}
	dishHeight = row
	dishWidth += 1

	tiltDishNorth(dish)
	// printDish(dish)
	fmt.Println(scoreDish(dish))

	// Part two - complete first cycle
	tiltDishWest(dish)
	tiltDishSouth(dish)
	tiltDishEast(dish)
	// fmt.Println("After 1 cycle:")
	// printDish(dish)

	// Assume at some point we will loop
	// dish state => iteration it was seen
	seenDishes := make(map[string]int)
	scores := []int{0} // Ignore "0th" iteration
	iterations := 1
	seenDishes[dishState(dish)] = iterations
	scores = append(scores, scoreDish(dish))

	target := 1000000000
	for {
		cycleDish(dish)
		iterations += 1
		scores = append(scores, scoreDish(dish))

		stateStr := dishState(dish)
		if lastSeen, present := seenDishes[stateStr]; present {
			loopLen := iterations - lastSeen
			fmt.Println("\tFound loop of len ", loopLen, " after ", iterations, " iterations")

			target -= lastSeen
			target %= loopLen
			fmt.Println(scores[lastSeen+target])

			break
		} else {
			seenDishes[stateStr] = iterations
		}
	}
}

func cycleDish(dish map[Point]rune) {
	tiltDishNorth(dish)
	tiltDishWest(dish)
	tiltDishSouth(dish)
	tiltDishEast(dish)
}

func tiltDishNorth(dish map[Point]rune) {
	// Updates the map

	// Tilt north - everything rolls upwards
	// Do one column at a time
	for col := 0; col < dishWidth; col++ {
		// Go bottom to top
		firstFreeSpace := 0
		for row := 0; row < dishHeight; row++ {
			if r, present := dish[Point{row, col}]; present {
				if r == '#' {
					firstFreeSpace = row + 1
				} else if r == 'O' {
					delete(dish, Point{row, col})
					dish[Point{firstFreeSpace, col}] = 'O'
					firstFreeSpace += 1
				} else {
					panic("Unexpected rune: " + string(r))
				}
			}
		}
	}
}

func tiltDishSouth(dish map[Point]rune) {
	// Updates the map
	for col := 0; col < dishWidth; col++ {
		// Go top to bottom
		firstFreeSpace := dishHeight - 1
		for row := dishHeight - 1; row >= 0; row-- {
			if r, present := dish[Point{row, col}]; present {
				if r == '#' {
					firstFreeSpace = row - 1
				} else if r == 'O' {
					delete(dish, Point{row, col})
					dish[Point{firstFreeSpace, col}] = 'O'
					firstFreeSpace -= 1
				} else {
					panic("Unexpected rune: " + string(r))
				}
			}
		}
	}
}

func tiltDishWest(dish map[Point]rune) {
	// Updates the map
	for row := 0; row < dishHeight; row++ {
		// Go left to right
		firstFreeSpace := 0
		for col := 0; col < dishWidth; col++ {
			if r, present := dish[Point{row, col}]; present {
				if r == '#' {
					firstFreeSpace = col + 1
				} else if r == 'O' {
					delete(dish, Point{row, col})
					dish[Point{row, firstFreeSpace}] = 'O'
					firstFreeSpace += 1
				} else {
					panic("Unexpected rune: " + string(r))
				}
			}
		}
	}
}

func tiltDishEast(dish map[Point]rune) {
	// Updates the map
	for row := 0; row < dishHeight; row++ {
		// Go right to left
		firstFreeSpace := dishWidth - 1
		for col := dishWidth - 1; col >= 0; col-- {
			if r, present := dish[Point{row, col}]; present {
				if r == '#' {
					firstFreeSpace = col - 1
				} else if r == 'O' {
					delete(dish, Point{row, col})
					dish[Point{row, firstFreeSpace}] = 'O'
					firstFreeSpace -= 1
				} else {
					panic("Unexpected rune: " + string(r))
				}
			}
		}
	}
}

func scoreDish(dish map[Point]rune) int {
	score := 0
	for k, v := range dish {
		if v == 'O' {
			score += dishHeight - k.row
		}
	}
	return score
}

func printDish(dish map[Point]rune) {
	for row := 0; row < dishHeight; row++ {
		line := ""
		for col := 0; col < dishWidth; col++ {
			if r, present := dish[Point{row, col}]; present {
				line += string(r)
			} else {
				line += "."
			}
		}
		fmt.Println(line)
	}
}

func dishState(dish map[Point]rune) string {
	// Convert the state into a key that can be used in a map
	// Ignore hard rocks (#) as they never move
	oPoints := []string{}
	for k, v := range dish {
		if v == 'O' {
			oPoints = append(oPoints, fmt.Sprintf("%d-%d", k.row, k.col))
		}
	}
	slices.Sort(oPoints)

	return strings.Join(oPoints, ",")
}
