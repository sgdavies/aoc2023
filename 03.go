package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strconv"
)

type Point struct {
	row int
	col int
}

func (p Point) Adjacents() []Point {
	return []Point{
		Point{p.row-1, p.col-1},
		Point{p.row-1, p.col},
		Point{p.row-1, p.col+1},
		Point{p.row, p.col-1},
		Point{p.row, p.col+1},
		Point{p.row+1, p.col-1},
		Point{p.row+1, p.col},
		Point{p.row+1, p.col+1},
	}
}

type number struct {
	id int
	value int
}

func day03() {
	file, err := os.Open("data/03.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// . . 1 3  } maps to numbers {(0,2)=13, (0,3)=13}
	// . . + .  }         symbols {(1,2)='+'}
	numbers := make(map[Point]number)
	symbols := make(map[Point]rune)

	nextId := 0

	for row, line := range lines {
		var currentNumber *number = nil
		for col, c := range line {
			if c == '.' {
				currentNumber = nil
			} else if '0'<=c && c<='9' {
				if currentNumber == nil {
					num := string(c)
					for _, c2 := range(line[col+1:]) {
						if !('0'<=c2 && c2<='9') {
							break
						} else {
							num += string(c2)
						}
					}
					value, _ := strconv.Atoi(num)
					currentNumber = &number{nextId, value}
					nextId ++
				}

				numbers[Point{row,col}] = *currentNumber
			} else {
				currentNumber = nil
				symbols[Point{row,col}] = c
			}
		}
	}
	// fmt.Println(fmt.Sprint(len(numbers)))
	// fmt.Println(fmt.Sprint(len(symbols)))

	// Okay, we have the data. Now find the numbers that are adjacent to symbols.
	partOne := 0
	partTwo := 0
	for symPoint, sym := range symbols {
		adjs := make(map[int]number) // prevent double counting, and use for part two
		for _, adj := range symPoint.Adjacents() {
			if num, exists := numbers[adj]; exists {
				if _, seen := adjs[num.id]; !seen {
					partOne += num.value
					adjs[num.id] = num
				}
			}
		}
		// Part two
		if sym == '*'  && len(adjs) == 2 {
			gearRatio := 1
			for _, v := range adjs {
				gearRatio *= v.value
			}
			partTwo += gearRatio
		}
	}

	fmt.Println(partOne)
	fmt.Println(partTwo)
}