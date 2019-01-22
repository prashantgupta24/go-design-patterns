package main

import "fmt"

//INTERFACES

type Sound interface {
	MakeSound()
}

type Action interface {
	DoAction()
}

//IMPLEMENTATIONS
type Quack struct{}

func (q *Quack) MakeSound() {
	fmt.Println("Quack!")
}

type NoQuack struct{}

func (nq *NoQuack) MakeSound() {
	fmt.Println("No sound!")
}

type Fly struct{}

func (f *Fly) DoAction() {
	fmt.Println("Flying!")
}

type NoFly struct{}

func (nf *NoFly) DoAction() {
	fmt.Println("Can't fly!")
}

//Structs
type Duck struct {
	Sound
	Action
}

func main() {
	// d := &DuckType1{}
	// d.DoAction()

	q := &NoQuack{}
	f := &Fly{}
	d1 := Duck{
		Sound:  q,
		Action: f,
	}
	d1.MakeSound()
	d1.DoAction()
}
