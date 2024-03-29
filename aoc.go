package main

import (
	"fmt"
	"os"
	"strconv"
)

var days = [25]func(){
	day01,
	day02,
	day03,
	day04,
	day05,
	day06,
	day07,
	day08,
	day09,
	day10,
	day11,
	day12,
	day13,
	day14,
	day15,
	day16,
	day17,
	day18,
	day19,
	day20,
	day21,
	day22,
	day23,
	day24,
	day25,
}

func main() {
	fmt.Println("Ho Ho Ho, Advent of Code!")
	args := os.Args
	if len(args) == 1 {
		// Run all
		for i, f := range days {
			if f != nil {
				fmt.Println("Day " + fmt.Sprintf("%02d", i+1))
				f()
			}
		}
	} else {
		for _, d := range args[1:] {
			day, err := strconv.Atoi(d)
			if err != nil || day > len(days) || days[day-1] == nil {
				fmt.Println("Invalid day or not implemented: " + d)
				continue
			}

			fmt.Println("Day " + fmt.Sprintf("%02d", day))
			days[day-1]()
		}
	}
}
