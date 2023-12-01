package main

import (
	"fmt"
	"os"
	"strconv"
)

var days = [25]func(){
	day01,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
	nil,
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
			if err != nil {
				err = fmt.Errorf("Invalid day '"+d+"': ", err)
				panic(err)
			}
			if day > len(days) || days[day-1] == nil {
				panic("Day " + fmt.Sprint(day) + " not implemented")
			}

			days[day-1]()
		}
	}
}
