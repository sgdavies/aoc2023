package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func day12() {
	file, err := os.Open("data/12.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	partOne := 0
	partTwo := 0
	progress := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		groupStr, targetStr := splitSprings(line)
		groups, target := parseSprings(groupStr, targetStr)
		asum := countWays(groups, target)
		// fmt.Println(line, asum)
		partOne += asum

		// Part two
		twoGroupStr := groupStr
		twoTargetStr := targetStr
		groups[len(groups)-1] = groups[len(groups)-1] + "?"
		for i := 0; i < 4; i++ {
			twoGroupStr = twoGroupStr + "?" + groupStr
			twoTargetStr = twoTargetStr + "," + targetStr
		}

		twoGroups, twoTarget := parseSprings(twoGroupStr, twoTargetStr)
		bsum := countWays(twoGroups, twoTarget)
		partTwo += bsum

		progress += 1
		if progress%10 == 0 {
			if progress%100 == 0 {
				fmt.Println(progress)
			} else {
				fmt.Print(progress)
			}
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()

	fmt.Println(partOne)
	fmt.Println(partTwo)
}

func splitSprings(line string) (string, string) {
	parts := strings.Split(line, " ")
	return parts[0], parts[1]
}

func parseSprings(groupStr, targetStr string) ([]string, []int) {
	// .??..??...?##. 1,1,3 ==> ["??","??","?##"]  [1,1,3]
	groupsWithDots := strings.Split(groupStr, ".")
	groups := []string{}
	for _, g := range groupsWithDots {
		if len(g) > 0 {
			groups = append(groups, g)
		}
	}
	target := []int{}
	tParts := strings.Split(targetStr, ",")
	for _, tp := range tParts {
		v, _ := strconv.Atoi(tp)
		target = append(target, v)
	}

	return groups, target
}

func countWays(groups []string, target []int) int {
	totalWays := 0
	for ig, g := range groups {
		for _, remnant := range springMatch(g, target[0]) {
			newGroups := groups[ig+1:]
			if len(remnant) > 0 {
				newGroups = append([]string{remnant}, newGroups...)
			}

			if len(target) == 1 {
				// This was the last target. This is a solution
				// if there are no #s left in the remaining groups.
				if !containsDefiniteSprings(newGroups) {
					totalWays += 1
				}
			} else if len(newGroups) > 0 {
				totalWays += countWays(newGroups, target[1:])
			}
		}

		if containsDefiniteSprings([]string{g}) {
			// Can't skip past this - it must be consumed
			break
		}
	}

	return totalWays
}

func springMatch(group string, tlen int) []string {
	ret := []string{}
	for i := range group {
		rGroup := group[i:]
		if tlen > len(rGroup) {
			break
		}

		// It fits. There must be a gap after it (not #)
		if tlen == len(rGroup) {
			ret = append(ret, "")
		} else {
			// There must be a gap after the group
			if rGroup[tlen] != '#' {
				ret = append(ret, rGroup[tlen+1:])
			}
		}

		// We must consume any definite springs, so we can't skip past them
		if rGroup[0] == '#' {
			break
		}
	}

	return ret
}

func containsDefiniteSprings(groups []string) bool {
	for _, group := range groups {
		if strings.Contains(group, "#") {
			return true
		}
	}
	return false
}
