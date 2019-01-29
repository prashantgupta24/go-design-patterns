package main

import "fmt"

type tt struct {
	s string
}

func (t *tt) init() {
	fmt.Println("init")
	t.s = "abc"
}
func main() {
	t := &tt{}
	t.init()
	fmt.Println("main")
}
