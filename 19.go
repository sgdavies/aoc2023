package main

import (
	"fmt"
	"strconv"
	"strings"
)

type part struct {
	x int
	m int
	a int
	s int
}

func day19() {
	lines := LinesFromFile("data/19.txt")
	parts := []part{}
	// accepted := []part{}
	rules := make(map[string][]string)

	i := 0
	for {
		line := lines[i]
		i += 1
		if line == "" {
			break
		}

		// px{a<2006:qkq,m>2090:A,rfg}
		bits := strings.Split(line, "{")
		rule := bits[1]
		rule = rule[:len(rule)-1]
		rulez := strings.Split(rule, ",")
		rules[bits[0]] = rulez
	}

	for _, line := range lines[i:] {
		parts = append(parts, parsePart(line))
	}

	acceptedSum := 0
	for _, p := range parts {
		val := processPart(p, rules)
		// fmt.Println(val)
		acceptedSum += val
	}

	fmt.Println(acceptedSum)

	// Part two
	start := part{1, 1, 1, 1}
	end := part{4000, 4000, 4000, 4000}
	fmt.Println(solve19recurse(start, end, "in", 0, rules))
}

func parsePart(line string) part {
	// {x=787,m=2655,a=1222,s=2876}
	bits := strings.Split(line, ",")
	x, _ := strconv.Atoi(bits[0][3:])
	m, _ := strconv.Atoi(bits[1][2:])
	a, _ := strconv.Atoi(bits[2][2:])
	s, _ := strconv.Atoi(bits[3][2 : len(bits[3])-1])

	return part{x, m, a, s}
}

func processPart(p part, rules map[string][]string) int {
	// a<2006:qkq,m>2090:A,rfg
	ruleName := "in"
	// fmt.Println(p)
	for {
		// fmt.Print(ruleName + " ")
		rule := rules[ruleName]
		for _, r := range rule {
			if r == "R" {
				return 0
			}

			if r == "A" {
				return p.x + p.m + p.a + p.s
			}

			// End of rule - it's the name of the next rule to apply
			if !strings.Contains(r, ":") {
				ruleName = r
				break
			}

			// Parse and apply the rule
			var val int
			switch r[0] {
			case 'x':
				{
					val = p.x
				}
			case 'm':
				{
					val = p.m
				}
			case 'a':
				{
					val = p.a
				}
			case 's':
				{
					val = p.s
				}
			default:
				{
					panic("Invalid [xmas] for rule: " + r)
				}
			}

			// a<2006:qkq or m>2090:A
			op := r[1]
			bits := strings.Split(r[2:], ":")
			cmp, _ := strconv.Atoi(bits[0])
			nextR := bits[1]
			applies := (op == '<' && val < cmp) || (op == '>' && val > cmp)
			if applies {
				if nextR == "R" {
					return 0
				}

				if nextR == "A" {
					return p.x + p.m + p.a + p.s
				}

				ruleName = nextR
				break
			}
		}
	}
}

func solve19recurse(start part, end part, ruleName string, ruleNumber int, rules map[string][]string) int {
	rule := rules[ruleName]
	r := rule[ruleNumber]

	if r == "A" {
		return combo19(start, end)
	}

	if r == "R" {
		return 0
	}

	if !strings.Contains(r, ":") {
		return solve19recurse(start, end, r, 0, rules)
	}

	// a<2006:qkq or m>2090:A
	// split into new parts covering each case, and recurse if necessary
	xmas := r[0]
	op := r[1]
	bits := strings.Split(r[2:], ":")
	cmp, _ := strconv.Atoi(bits[0])
	nextR := bits[1]

	var startA, endA = start, end
	var startB, endB = start, end

	var a, b int
	if op == '<' {
		a = cmp - 1
		b = cmp
	} else {
		a = cmp
		b = cmp + 1
	}

	switch xmas {
	case 'x':
		{
			endA.x = a
			startB.x = b
		}
	case 'm':
		{
			endA.m = a
			startB.m = b
		}
	case 'a':
		{
			endA.a = a
			startB.a = b
		}
	case 's':
		{
			endA.s = a
			startB.s = b
		}
	}

	// The one that matches the rule recurses with the new rule number (or combo);
	// the other recurses with the next rule index
	valA, valB := 0, 0
	if op == '<' {
		// A covers the matching case, B covers the non-matching case
		if startA.x <= endA.x && startA.m <= endA.m && startA.a <= endA.a && startA.s <= endA.s {
			if nextR == "A" {
				valA = combo19(startA, endA)
			} else if nextR == "R" {
				valA = 0
			} else {
				valA = solve19recurse(startA, endA, nextR, 0, rules)
			}
		}

		if startB.x <= endB.x && startB.m <= endB.m && startB.a <= endB.a && startB.s <= endB.s {
			valB = solve19recurse(startB, endB, ruleName, ruleNumber+1, rules)
		}
	} else if op == '>' {
		// A covers non-matching case, B covers the matching case
		if startA.x <= endA.x && startA.m <= endA.m && startA.a <= endA.a && startA.s <= endA.s {
			valA = solve19recurse(startA, endA, ruleName, ruleNumber+1, rules)
		}

		if startB.x <= endB.x && startB.m <= endB.m && startB.a <= endB.a && startB.s <= endB.s {
			if nextR == "A" {
				valB = combo19(startB, endB)
			} else if nextR == "R" {
				valB = 0
			} else {
				valB = solve19recurse(startB, endB, nextR, 0, rules)
			}
		}
	}

	return valA + valB
}

func combo19(start part, end part) int {
	val := (end.x - start.x + 1) * (end.m - start.m + 1) * (end.a - start.a + 1) * (end.s - start.s + 1)
	// fmt.Println(val)
	return val
}
