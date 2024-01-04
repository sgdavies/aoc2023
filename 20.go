package main

import (
	"fmt"
	"os"
	"strings"
)

type pulse int

const (
	loPulse pulse = iota
	highPulse
)

type signal struct {
	value pulse
	from  string
	to    string
}

type flipflop struct {
	state int // 0 = off, 1 = on
	dests []string
}

func (f flipflop) Dests() []string {
	return f.dests
}

type conjunction struct {
	inputs map[string]pulse
	dests  []string
}

func (c conjunction) Dests() []string {
	return c.dests
}

type plain struct {
	dests []string
}

func (p plain) Dests() []string {
	return p.dests
}

type module interface {
	Dests() []string
}

func day20() {
	lines := lines20("real")
	n2n := make(map[string]module) // node name to node object

	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		name := parts[0]
		dests := strings.Split(parts[1], ", ")
		switch name[0] {
		case '%':
			{
				n2n[name[1:]] = &flipflop{0, dests}
			}
		case '&':
			{
				n2n[name[1:]] = &conjunction{map[string]pulse{}, dests}
			}
		default:
			{
				n2n[name] = &plain{dests}
			}
		}
	}

	// add conjunction inputs & collect any missing (plain) nodes
	for k, v := range n2n {
		for _, d := range v.Dests() {
			dd, present := n2n[d]
			if !present {
				n2n[d] = &plain{nil}
			}
			dd = n2n[d]
			dc, ok := dd.(*conjunction)
			if ok {
				dc.inputs[k] = loPulse
			}
		}
	}

	// Solve the puzzle
	lo, hi := 0, 0
	presses := 0
	bbInputs := make(map[string]int)
	for {
		if presses == 1000 {
			fmt.Println(lo * hi)
			// break
		}
		presses++
		rxLoCount := 0
		signals := []signal{{loPulse, "button", "broadcaster"}}
		for len(signals) > 0 {
			sig := signals[0]
			signals = signals[1:]
			if sig.value == loPulse {
				lo++
			} else if sig.value == highPulse {
				hi++
			} else {
				panic("Invalid signal value")
			}

			if sig.to == "rx" && sig.value == loPulse {
				rxLoCount++
			}

			receiver, ok := n2n[sig.to]
			if !ok {
				panic("No record of node: " + sig.to)
			}
			var p pulse // pulse to send
			dests := receiver.Dests()
			switch receiver.(type) {
			case *flipflop:
				{
					if sig.value == loPulse {
						rff := receiver.(*flipflop)
						if rff.state == 1 {
							rff.state = 0
							p = loPulse
						} else {
							rff.state = 1
							p = highPulse
						}
					} else {
						dests = nil
					}
				}
			case *conjunction:
				{
					rc := receiver.(*conjunction)
					rc.inputs[sig.from] = sig.value
					il, ih := 0, 0
					for _, vp := range rc.inputs {
						if vp == loPulse {
							il++
						} else {
							ih++
						}
					}
					if il == 0 {
						p = loPulse
					} else {
						p = highPulse
					}
					if strings.Contains("ct kp ks xc", sig.to) && p == highPulse {
						if _, pres := bbInputs[sig.to]; !pres {
							bbInputs[sig.to] = presses
						}
						if len(bbInputs) == 4 {
							// fmt.Println(bbInputs)
							partTwo := 1
							for _, v := range bbInputs {
								partTwo *= v
							}
							fmt.Println(partTwo)
							os.Exit(0)
						}
					}
				}
			case *plain:
				{
					p = sig.value
				}
			default:
				{
					panic("Didn't match a type")
				}
			}
			for _, d := range dests {
				signals = append(signals, signal{p, sig.to, d})
			}
		}

		if rxLoCount > 0 {
			fmt.Print(rxLoCount)
			if rxLoCount == 1 {
				fmt.Println(presses)
				os.Exit(0)
			}
		}
	}
	// fmt.Println( /*lo, hi,*/ lo * hi)
}

func lines20(source string) []string {
	switch source {
	case "ex1":
		{
			ex1 := `broadcaster -> a, b, c
%a -> b
%b -> c
%c -> inv
&inv -> a`
			return strings.Split(ex1, "\n")
		}
	case "ex2":
		{
			ex2 := `broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output`
			return strings.Split(ex2, "\n")
		}
	case "real":
		{
			return LinesFromFile("data/20.txt")
		}
	default:
		panic("Unexpected source name: " + source)
	}
}
