package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func day02() {
	file, err := os.Open("data/02.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	possibleSum := 0
	powerSum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		gameId, possible, power := evaluateBag(line)
		if possible {
			// possibles = append(possibles, gameId)
			possibleSum += gameId
		}
		powerSum += power
	}

	fmt.Println(possibleSum)
	fmt.Println(powerSum)
}

func evaluateBag(line string) (int, bool, int) {
	// fmt.Println(line)
	possible := true
	mins := [3]int{0, 0, 0} // r,g,b

	parts := strings.Split(line, ": ")
	gameId, _ := strconv.Atoi(parts[0][5:]) // strip "Game "
	rounds := strings.Split(parts[1], "; ")
	for _, round := range rounds {
		balls := strings.Split(round, ", ")
		for _, ball := range balls {
			bits := strings.Split(ball, " ")
			number, _ := strconv.Atoi(bits[0])
			color := bits[1]
			switch color {
			case "red":
				{
					if number > 12 {
						possible = false
					}

					if number > mins[0] {
						mins[0] = number
					}
				}
			case "green":
				{
					if number > 13 {
						possible = false
					}

					if number > mins[1] {
						mins[1] = number
					}
				}
			case "blue":
				{
					if number > 14 {
						possible = false
					}

					if number > mins[2] {
						mins[2] = number
					}
				}
			default:
				{
					panic("Unexpected color " + color + " in line '" + line + "'")
				}
			}
		}
	}

	return gameId, possible, mins[0] * mins[1] * mins[2]
}
