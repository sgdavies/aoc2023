package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func day13() {
	file, err := os.Open("data/13.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	blocks := [][]string{}
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			lines = append(lines, line)
		} else {
			// We have a complete block
			blocks = append(blocks, lines)
			lines = []string{} // Reset for next block
		}
	}
	// ... and the final block (not caught above by the "" line)
	blocks = append(blocks, lines)

	partOne := 0
	partTwo := 0
	for _, block := range blocks {
		hcols, vrows := solveMirrorBlock(block)
		partOne += hcols + 100*vrows

		// Part two
		// Need to find a different reflection line - but the old one may still be valid
		foundDifferent := false
		for _, mblock := range modifiedBlocks(block) {
			// Assume single solution ... ?
			mhcols := findMirror(transpose(mblock), hcols-1)
			mvrows := findMirror(mblock, vrows-1)
			if mhcols > 0 && mhcols != hcols {
				partTwo += mhcols
				foundDifferent = true
				break
			} else if mvrows > 0 && mvrows != vrows {
				partTwo += 100 * mvrows
				foundDifferent = true
				break
			}
		}
		if !foundDifferent {
			panic("Couldn't find different reflection (" + fmt.Sprint(hcols) + ", " + fmt.Sprint(vrows) + ") :\n" + strings.Join(block, "\n"))
		}
	}
	fmt.Println(partOne)
	fmt.Println(partTwo)
}

func solveMirrorBlock(lines []string) (int, int) {
	hcols := findMirror(transpose(lines), -1)
	vrows := findMirror(lines, -1)
	if hcols > 0 && vrows > 0 {
		panic("Can't have two mirrors!")
	} else if hcols == 0 && vrows == 0 {
		panic("Must have at least one mirror!")
	}
	return hcols, vrows
}

func transpose(lines []string) []string {
	tlines := make([]string, len(lines[0]))
	for _, line := range lines {
		for col, r := range line {
			tlines[col] = tlines[col] + string(r)
		}
	}

	return tlines
}

func findMirror(lines []string, partOne int) int {
	// partOne = partOne answer. Ignore this one. Set to -1 to skip.

	// If there are multiple lines of reflection - pick the one with the longest
	// reflection length (underspecified in puzzle)
	reflectPoint := -1
	longestLen := 0
	for row, line := range lines[:len(lines)-1] {
		if row == partOne {
			continue
		}
		if line == lines[row+1] {
			// if reflectPoint != -1 {
			// 	panic("Two lines of reflection: " + fmt.Sprint(reflectPoint) + " / " + fmt.Sprint(row) + "\n" + strings.Join(lines, "\n"))
			// }

			perfectMirror := true
			length := 1
			for row-length >= 0 && row+1+length < len(lines) {
				if lines[row-length] != lines[row+1+length] {
					perfectMirror = false
					break
				}
				length += 1
			}
			if perfectMirror && length > longestLen {
				reflectPoint = row
				longestLen = length
			}
		}
	}
	return reflectPoint + 1
}

func modifiedBlocks(block []string) [][]string {
	// Generate every possible un-smudged block
	mblocks := [][]string{}
	for row, line := range block {
		for col, r := range line {
			var newr string
			if r == '#' {
				newr = "."
			} else {
				newr = "#"
			}

			newline := line[:col] + newr + line[col+1:]
			mblock := make([]string, len(block)) // Copy - don't append onto original block
			copy(mblock, block)
			mblock = append(mblock[:row], newline)
			mblock = append(mblock, block[row+1:]...)
			mblocks = append(mblocks, mblock)
		}
	}
	return mblocks
}
