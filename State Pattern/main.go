package main

import "fmt"

type state string

const (
	noQuarter  state = "NO_QUARTER"
	hasQuarter state = "HAS_QUARTER"
)

type machine struct {
	state state
}

func (m *machine) insertQuarter() {
	fmt.Println("inserting quarter ...")
	if m.state == noQuarter {
		m.state = hasQuarter
	} else if m.state == hasQuarter {
		fmt.Println("machine has quarter already, please wait ...")
	}
}
func main() {
	fmt.Println("starting machine")

	m := &machine{state: noQuarter}
	m.insertQuarter()
	m.insertQuarter()
}
