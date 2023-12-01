package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var re_day01_nums_only = regexp.MustCompile(`[^0-9]`)

func day01() {
	file, err := os.Open("data/01.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sum1 := 0
	sum2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sum1 += value(line)
		two := 10*value2first(line) + value2last(line)
		sum2 += two
	}

	fmt.Println(sum1)
	fmt.Println(sum2)
}

func value(line string) int {
	nums := re_day01_nums_only.ReplaceAllString(line, "")
	n1, _ := strconv.Atoi(string(nums[0]))
	n2, _ := strconv.Atoi(string(nums[len(nums)-1]))
	return n1*10 + n2
}

func value2first(line string) int {
	first := -1
	for i, v := range line {
		if v >= '0' && v <= '9' {
			first, _ = strconv.Atoi(string(v))
			break
		} else if strings.HasPrefix(line[i:], "zero") {
			first = 0
			break
		} else if strings.HasPrefix(line[i:], "one") {
			first = 1
			break
		} else if strings.HasPrefix(line[i:], "two") {
			first = 2
			break
		} else if strings.HasPrefix(line[i:], "three") {
			first = 3
			break
		} else if strings.HasPrefix(line[i:], "four") {
			first = 4
			break
		} else if strings.HasPrefix(line[i:], "five") {
			first = 5
			break
		} else if strings.HasPrefix(line[i:], "six") {
			first = 6
			break
		} else if strings.HasPrefix(line[i:], "seven") {
			first = 7
			break
		} else if strings.HasPrefix(line[i:], "eight") {
			first = 8
			break
		} else if strings.HasPrefix(line[i:], "nine") {
			first = 9
			break
		}
	}

	if first < 0 {
		log.Fatal("Couldn't find first value in line: " + line)
	}

	// log.Println("Line: " + line + " -> " + fmt.Sprint(first))
	return first
}

func value2last(line string) int {
	last := -1
	for j, _ := range line {
		i := len(line) - 1 - j
		v := line[i]
		if v >= '0' && v <= '9' {
			last, _ = strconv.Atoi(string(v))
			break
		} else if strings.HasPrefix(line[i:], "zero") {
			last = 0
			break
		} else if strings.HasPrefix(line[i:], "one") {
			last = 1
			break
		} else if strings.HasPrefix(line[i:], "two") {
			last = 2
			break
		} else if strings.HasPrefix(line[i:], "three") {
			last = 3
			break
		} else if strings.HasPrefix(line[i:], "four") {
			last = 4
			break
		} else if strings.HasPrefix(line[i:], "five") {
			last = 5
			break
		} else if strings.HasPrefix(line[i:], "six") {
			last = 6
			break
		} else if strings.HasPrefix(line[i:], "seven") {
			last = 7
			break
		} else if strings.HasPrefix(line[i:], "eight") {
			last = 8
			break
		} else if strings.HasPrefix(line[i:], "nine") {
			last = 9
			break
		}
	}

	if last < 0 {
		log.Fatal("Couldn't find last value in line: " + line)
	}

	// log.Println("Line: " + line + " -> " + fmt.Sprint(last))
	return last
}
