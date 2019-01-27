package main

import "fmt"

type beverage interface {
	GetCost() float32
	GetDescription() string
}

type beverageStruct struct {
	description string
	cost        float32
}

func (b *beverageStruct) GetCost() float32 {
	return b.cost
}

func (b *beverageStruct) GetDescription() string {
	return b.description
}

type mocha struct {
	*beverageStruct
}

func (m *mocha) GetCost() float32 {
	return m.cost + float32(0.2)
}

func (m *mocha) GetDescription() string {
	return m.description + " + mocha"
}

func mochaDecorator(b beverage) beverage {
	beverageWithMocha := &beverageStruct{}
	beverageWithMocha.description = b.GetDescription() + " mocha"
	beverageWithMocha.cost = b.GetCost() + float32(0.2)
	return beverageWithMocha
}

func display(b beverage) {
	fmt.Printf("Type : %v and cost : %v\n", b.GetDescription(), b.GetCost())
	//fmt.Println("result ", float32(1.1)+float32(1.2))

}
func main() {

	b := &beverageStruct{
		description: "Expresso",
		cost:        3,
	}
	display(b)

	bWithMochaStruct := &mocha{
		beverageStruct: b,
	}
	display(bWithMochaStruct)

	bWithMochaFunc := mochaDecorator(b)
	display(bWithMochaFunc)
}
