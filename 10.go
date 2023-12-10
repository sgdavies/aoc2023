package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type direction int

const (
	above direction = iota
	left
	below
	right
)

func day10() {
	file, err := os.Open("data/10.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var start Point
	pieces := make(map[Point]rune)
	for row, line := range lines {
		for col, r := range line {
			if r == '.' {
				continue
			}
			p := Point{row, col}
			if r == 'S' {
				start = p
			}
			pieces[p] = r
		}
	}

	loop := []Point{start}
	var nextPiece Point
	var nextDirection direction // So we know which way to go from next piece (not backwards)
	// Find a starting piece
	pAbove := Point{start.row - 1, start.col}
	pRight := Point{start.row, start.col + 1}
	pBelow := Point{start.row + 1, start.col}
	// There are 2 connecting pieces, so need to check max 3 of 4
	if piece, ok := pieces[pAbove]; ok && slices.Contains([]rune{'|', '7', 'F'}, piece) {
		nextPiece = pAbove
		nextDirection = above
	} else if piece, ok := pieces[pRight]; ok && slices.Contains([]rune{'-', '7', 'J'}, piece) {
		nextPiece = pRight
		nextDirection = right
	} else if piece, ok := pieces[pBelow]; ok && slices.Contains([]rune{'|', 'J', 'L'}, piece) {
		nextPiece = pBelow
		nextDirection = below
	} else {
		panic("Can't find piece connected to start")
	}

	firstDirection := nextDirection // to identify S later
	for nextPiece != start {
		loop = append(loop, nextPiece)
		pipe, _ := pieces[nextPiece]
		switch pipe {
		case '|':
			{
				if nextDirection == above {
					nextPiece = Point{nextPiece.row - 1, nextPiece.col}
					nextDirection = above
				} else if nextDirection == below {
					nextPiece = Point{nextPiece.row + 1, nextPiece.col}
					nextDirection = below
				} else {
					panic("Can't get here from left or right")
				}
			}
		case '-':
			{
				if nextDirection == left {
					nextPiece = Point{nextPiece.row, nextPiece.col - 1}
					nextDirection = left
				} else if nextDirection == right {
					nextPiece = Point{nextPiece.row, nextPiece.col + 1}
					nextDirection = right
				} else {
					panic("Can't get here from above or below")
				}
			}
		case 'J':
			{
				if nextDirection == right {
					nextPiece = Point{nextPiece.row - 1, nextPiece.col}
					nextDirection = above
				} else if nextDirection == below {
					nextPiece = Point{nextPiece.row, nextPiece.col - 1}
					nextDirection = left
				} else {
					panic("Can't get here from right or below")
				}
			}
		case 'L':
			{
				if nextDirection == left {
					nextPiece = Point{nextPiece.row - 1, nextPiece.col}
					nextDirection = above
				} else if nextDirection == below {
					nextPiece = Point{nextPiece.row, nextPiece.col + 1}
					nextDirection = right
				} else {
					panic("Can't get here from left or below")
				}
			}
		case '7':
			{
				if nextDirection == right {
					nextPiece = Point{nextPiece.row + 1, nextPiece.col}
					nextDirection = below
				} else if nextDirection == above {
					nextPiece = Point{nextPiece.row, nextPiece.col - 1}
					nextDirection = left
				} else {
					panic("Can't get here from right or above")
				}
			}
		case 'F':
			{
				if nextDirection == left {
					nextPiece = Point{nextPiece.row + 1, nextPiece.col}
					nextDirection = below
				} else if nextDirection == above {
					nextPiece = Point{nextPiece.row, nextPiece.col + 1}
					nextDirection = right
				} else {
					panic("Can't get here from above or left")
				}
			}

		default:
			{
				panic("Unexpected pipe: " + string(pipe))
			}
		}
	}

	fmt.Println(len(loop) / 2)

	// Part two
	// Plan:
	// Walk round the loop, 'touching' each adjacent point inside the loop
	// Then, flood-fill from each of these points.

	// First: find if our loop goes clockwise or anticlockwise
	// Start by redefining S to be the correct pipe type
	var s rune
	var dirSToPrev direction
	switch nextDirection {
	case above:
		dirSToPrev = below
	case below:
		dirSToPrev = above
	case left:
		dirSToPrev = right
	case right:
		dirSToPrev = left
	}
	// Sort
	if firstDirection < dirSToPrev {
		firstDirection, dirSToPrev = dirSToPrev, firstDirection
	}

	switch firstDirection {
	case above:
		{
			switch dirSToPrev {
			case left:
				s = 'J'
			case right:
				s = 'L'
			case below:
				s = '|'
			case above:
				panic("Can't both be above")
			}
		}
	case below:
		{
			switch dirSToPrev {
			case left:
				s = '7'
			case right:
				s = 'F'
			case above:
				s = '|'
			case below:
				panic("Can't both be below")
			}
		}
	case left:
		{
			switch dirSToPrev {
			case above:
				s = 'J'
			case right:
				s = '-'
			case below:
				s = '7'
			case left:
				panic("Can't both be left")
			}
		}
	case right:
		{
			switch dirSToPrev {
			case left:
				s = '-'
			case above:
				s = 'L'
			case below:
				s = 'F'
			case right:
				panic("Can't both be right")
			}
		}
	}

	pieces[start] = s

	// Walk round the loop, tagging the adjacent squares on left & right.
	// Keep track of turns so at the end we know which of left/right was on the inside
	leftSquares := make(map[Point]bool)
	rightSquares := make(map[Point]bool)
	turns := 0              // +ve = CW (turning right)
	facing := nextDirection // Facing into S from final piece
	for i, p := range loop {
		piece, _ := pieces[p]

		facing, turns = evaluatePipe(p, piece, facing, turns, leftSquares, rightSquares)

		turns += i
		turns -= i
	}

	var insideSquares map[Point]bool
	if turns == 4 {
		insideSquares = rightSquares
	} else if turns == -4 {
		insideSquares = leftSquares
	} else {
		panic("Invalid number of turns " + fmt.Sprint(turns))
	}

	// fmt.Println(len(insideSquares))
	// If pipe runs adjacent, we may have added in other piece of pipe - remove these
	for _, p := range loop {
		delete(insideSquares, p)
	}

	newSquares := make(map[Point]bool)
	for k := range insideSquares {
		newSquares[k] = true
	}
	seenSquares := insideSquares
	for _, p := range loop {
		seenSquares[p] = true
	}

	insideCount := 0
	for len(newSquares) > 0 {
		oldNewSquares := newSquares
		newSquares = make(map[Point]bool)
		for p, _ := range oldNewSquares {
			insideCount += 1
			for _, op := range []Point{{p.row - 1, p.col}, {p.row + 1, p.col}, {p.row, p.col - 1}, {p.row, p.col + 1}} {
				if _, seen := seenSquares[op]; !seen {
					newSquares[op] = true
					seenSquares[op] = true
				}
			}
		}
	}

	fmt.Println(insideCount)
}

func evaluatePipe(p Point, piece rune, facing direction, turns int, leftSquares map[Point]bool, rightSquares map[Point]bool) (direction, int) {
	switch piece {
	case '-':
		{
			if facing == right {
				leftSquares[Point{p.row - 1, p.col}] = true
				rightSquares[Point{p.row + 1, p.col}] = true
			} else if facing == left {
				leftSquares[Point{p.row + 1, p.col}] = true
				rightSquares[Point{p.row - 1, p.col}] = true
			} else {
				panic("Facing invalid direction for -")
			}
		}
	case '|':
		{
			if facing == above {
				leftSquares[Point{p.row, p.col - 1}] = true
				rightSquares[Point{p.row, p.col + 1}] = true
			} else if facing == below {
				leftSquares[Point{p.row, p.col + 1}] = true
				rightSquares[Point{p.row, p.col - 1}] = true
			} else {
				panic("Facing invalid direction for |")
			}
		}
	case 'F':
		{
			if facing == above {
				turns += 1
				facing = right

				leftSquares[Point{p.row, p.col - 1}] = true
				leftSquares[Point{p.row - 1, p.col - 1}] = true
				leftSquares[Point{p.row - 1, p.col}] = true
			} else if facing == left {
				turns -= 1
				facing = below

				rightSquares[Point{p.row - 1, p.col}] = true
				rightSquares[Point{p.row - 1, p.col - 1}] = true
				rightSquares[Point{p.row, p.col - 1}] = true
			} else {
				panic("Facing invalid direction for F")
			}
		}
	case 'J':
		{
			if facing == below {
				turns += 1
				facing = left

				leftSquares[Point{p.row, p.col + 1}] = true
				leftSquares[Point{p.row + 1, p.col + 1}] = true
				leftSquares[Point{p.row + 1, p.col}] = true
			} else if facing == right {
				turns -= 1
				facing = above

				rightSquares[Point{p.row + 1, p.col}] = true
				rightSquares[Point{p.row + 1, p.col + 1}] = true
				rightSquares[Point{p.row, p.col + 1}] = true
			} else {
				panic("Facing invalid direction for J")
			}
		}
	case 'L':
		{
			if facing == left {
				turns += 1
				facing = above

				leftSquares[Point{p.row + 1, p.col}] = true
				leftSquares[Point{p.row + 1, p.col - 1}] = true
				leftSquares[Point{p.row, p.col - 1}] = true
			} else if facing == below {
				turns -= 1
				facing = right

				rightSquares[Point{p.row, p.col - 1}] = true
				rightSquares[Point{p.row + 1, p.col - 1}] = true
				rightSquares[Point{p.row + 1, p.col}] = true
			} else {
				panic("Facing invalid direction for L")
			}
		}
	case '7':
		{
			if facing == right {
				turns += 1
				facing = below

				leftSquares[Point{p.row - 1, p.col}] = true
				leftSquares[Point{p.row - 1, p.col + 1}] = true
				leftSquares[Point{p.row, p.col + 1}] = true
			} else if facing == above {
				turns -= 1
				facing = left

				rightSquares[Point{p.row, p.col + 1}] = true
				rightSquares[Point{p.row - 1, p.col + 1}] = true
				rightSquares[Point{p.row - 1, p.col}] = true
			} else {
				panic("Facing invalid direction for 7")
			}
		}
	default:
		{
			panic("Unexpected piece!")
		}
	}

	return facing, turns
}
