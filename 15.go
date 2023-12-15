package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type lens struct {
	label string
	fl    int
}

func day15() {
	file, err := os.Open("data/15.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	// line = "rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7"
	steps := strings.Split(line, ",")
	sum := 0
	for _, step := range steps {
		sum += hash15(step)
	}

	fmt.Println(sum)

	// Part two
	hashMap := make([][]lens, 256)
	for _, step := range steps {
		if step[len(step)-1] == '-' {
			label := step[:len(step)-1]
			h := hash15(label)
			for i, l := range hashMap[h] {
				if l.label == label {
					lenses := hashMap[h]
					hashMap[h] = append(lenses[:i], lenses[i+1:]...)
					break
				}
			}
		} else {
			parts := strings.Split(step, "=")
			label := parts[0]
			value, _ := strconv.Atoi(parts[1])
			h := hash15(label)
			found := false
			for i, l := range hashMap[h] {
				if l.label == label {
					hashMap[h][i].fl = value
					found = true
					break
				}
			}

			if !found {
				hashMap[h] = append(hashMap[h], lens{label, value})
			}
		}
	}

	// fmt.Println(hashMap[0])
	// fmt.Println(hashMap[3])
	totalPower := 0
	for i, lenses := range hashMap {
		for j, l := range lenses {
			totalPower += (i + 1) * (j + 1) * l.fl
		}
	}
	fmt.Println(totalPower)
}

func hash15(step string) int {
	current := 0
	for _, r := range step {
		current += int(r)
		current *= 17
		current %= 256
	}
	return current
}
