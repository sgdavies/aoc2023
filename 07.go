package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type hand struct {
	cards string
	bid   int
	hand  int // 0 for high card, ..., 6 for 5-of-a-kind
}

func buildHand(line string) hand {
	var h hand
	parts := strings.Split(line, " ")
	h.cards = parts[0]
	h.bid, _ = strconv.Atoi(parts[1])
	h.hand = handValue(h.cards, false)
	return h
}

func handValue(cards string, jokers bool) int {
	cardCounts := make(map[rune]int)
	for _, r := range cards {
		if val, ok := cardCounts[r]; ok {
			cardCounts[r] = 1 + val
		} else {
			cardCounts[r] = 1
		}
	}

	// Number of singles, pairs, 3s, 4s, and 5s (ignore ix=0)
	groupCounts := []int{0, 0, 0, 0, 0, 0}
	for _, v := range cardCounts {
		groupCounts[v] += 1
	}
	// check
	sum := 0
	for i, count := range groupCounts {
		sum += i * count
	}
	if sum != 5 {
		panic("Should have 5 cards! " + cards)
	}

	// Account for jokers
	var njs int
	if val, ok := cardCounts['J']; ok {
		njs = val
	} else {
		njs = 0
	}
	if jokers && njs > 0 {
		switch njs {
		case 4:
			{
				groupCounts[4] = 0 // remove the jokers
				// add to sole other card to make 5-of-a-kind
				groupCounts[5] = 1
				groupCounts[1] = 0
			}
		case 3:
			{
				groupCounts[3] = 0 // remove the jokers
				if groupCounts[2] == 1 {
					groupCounts[2] = 0 // add to pair to make 5
					groupCounts[5] = 1
				} else {
					groupCounts[1] -= 1 // add to one of the singles to make 4
					groupCounts[4] = 1
				}
			}
		case 2:
			{
				groupCounts[2] -= 1
				if groupCounts[3] == 1 {
					groupCounts[3] = 0
					groupCounts[5] = 1
				} else if groupCounts[2] == 1 {
					groupCounts[2] = 0
					groupCounts[4] = 1
				} else {
					groupCounts[1] -= 1
					groupCounts[3] = 1
				}
			}
		case 1:
			{
				groupCounts[1] -= 1
				if groupCounts[4] == 1 {
					groupCounts[4] -= 1
					groupCounts[5] += 1
				} else if groupCounts[3] == 1 {
					groupCounts[3] = 0
					groupCounts[4] = 1
				} else if groupCounts[2] > 0 {
					groupCounts[2] -= 1
					groupCounts[3] = 1
				} else {
					groupCounts[1] -= 1
					groupCounts[2] = 1
				}
			}
		default:
			{
				// We already discounted 0, so this must be 5-of-a-kind
				if njs != 5 || groupCounts[5] != 1 {
					panic(cards)
				}
			}
		}
	}

	// Score the hand
	if groupCounts[5] == 1 {
		return 6
	} else if groupCounts[4] == 1 {
		return 5
	} else if groupCounts[3] == 1 && groupCounts[2] == 1 {
		return 4
	} else if groupCounts[3] == 1 {
		return 3
	} else if groupCounts[2] == 2 {
		return 2
	} else if groupCounts[2] == 1 {
		return 1
	} else {
		return 0
	}
}

func (a hand) cmp(b hand, jokers bool) int {
	if a.hand != b.hand {
		return a.hand - b.hand
	}

	for i := range a.cards {
		ar := a.cards[i]
		br := b.cards[i]
		if ar == br {
			continue
		} else {
			return cardVal(ar, jokers) - cardVal(br, jokers)
		}
	}

	return 0
}

func cardVal(c byte, jokers bool) int {
	if '2' <= c && c <= '9' {
		return int(c) - '0'
	} else {
		switch c {
		case 'T':
			return 10
		case 'J':
			{
				if jokers {
					return 1
				} else {
					return 11
				}
			}
		case 'Q':
			return 12
		case 'K':
			return 13
		case 'A':
			return 14
		default:
			{
				panic("Unexpected card " + string(c))
			}
		}
	}
}

func day07() {
	file, err := os.Open("data/07.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var hands []hand
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hands = append(hands, buildHand(scanner.Text()))
	}

	fmt.Println(day07score(hands, false))

	// Part two - rescore
	for i := range hands {
		hands[i].hand = handValue(hands[i].cards, true)
	}

	fmt.Println(day07score(hands, true))
}

func day07score(hands []hand, jokers bool) int {
	slices.SortFunc(hands, func(a, b hand) int { return a.cmp(b, jokers) })

	sum := 0
	for i, h := range hands {
		// fmt.Println(h.cards)
		sum += (i + 1) * h.bid
	}

	return sum
}
