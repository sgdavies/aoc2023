package main

import (
	"fmt"
	"math"
	"math/big"
)

var TEST = false

func day06() {
	var times, distances []int
	if TEST {
		times = []int{7, 15, 30}
		distances = []int{9, 40, 200}
	} else {
		times = []int{41, 66, 72, 66}
		distances = []int{244, 1047, 1228, 1040}
	}

	partOne := 1 // accumulator
	for i, time := range times {
		dist := distances[i]
		wins := 0
		for t := 1; t < time; t++ {
			if winRace(t, time, dist) {
				wins++
			}
		}
		partOne *= wins
	}

	fmt.Println(partOne)

	// Part two - solve it properly :)
	var time, distance int
	if TEST {
		time, distance = 71530, 940200
	} else {
		time, distance = 41667266, 244104712281040
	}
	// Math!
	// d = (t - p)*p > D  ==> -p^2 +t.p -D > 0
	// x = (-b +/- sqrt(b^2 - 4ac) / 2a) ; a=-1, b=t, c=-D
	// ==0 when p = (-t +/- sqrt(t^2 -4.D) ) / -2
	// ==0 when p = 0.5*(t +/- sqrt(t^2 -4.D))

	// Easy solution - normal ints (not good enough for real data)
	// lower := math.Ceil(0.5 * (float64(time) - math.Sqrt(float64(time*time-4*distance))))
	// upper := math.Floor(0.5 * (float64(time) + math.Sqrt(float64(time*time-4*distance))))
	// fmt.Println(upper - lower + 1)

	// Use big ints for real data
	lowerf := big.NewFloat(0.5 * (float64(time) - math.Sqrt(float64(time*time-4*distance))))
	upperf := big.NewFloat(0.5 * (float64(time) + math.Sqrt(float64(time*time-4*distance))))
	lower, _ := lowerf.Int(nil)
	if !lowerf.IsInt() {
		lower.Add(lower, big.NewInt(1)) // Int() rounds down - we want to round up to the first win
	}
	upper, _ := upperf.Int(nil) // Round down to get the last win
	// fmt.Println(lowerf, lower, upperf, upper)
	wins := upper.Sub(upper, lower)

	fmt.Println(wins.Add(wins, big.NewInt(1))) // +1 - count both ends of range
}

func winRace(press, time, dist int) bool {
	// if you hold button for 'press' will you go further than 'dist' in 'time'?
	speed := press
	travelled := speed * (time - press)
	return travelled > dist
}
