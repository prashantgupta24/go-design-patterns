package main

import "fmt"

type pizza struct {
	name string
	cost float32
}

func (p *pizza) Name() string {
	return p.name
}

func (p *pizza) Cost() float32 {
	return p.cost
}

type pizzaFactory interface {
	CreatePizza() *pizza
}

type nyPizzaFactory struct {
}

func (n *nyPizzaFactory) CreatePizza() *pizza {
	return &pizza{
		name: "NY pizza!",
		cost: 4,
	}
}

func chicagoPizzaCreator(cost float32) *pizza {
	// return func() *pizza {
	// 	return &pizza{
	// 		name: "Chicago pizza!",
	// 		cost: cost,
	// 	}
	// }
	return &pizza{
		name: "Chicago pizza!",
		cost: cost,
	}
}

type factoryFunc func(float32) *pizza
type factoryMap map[string]factoryFunc

func display(p *pizza) {
	fmt.Printf("Pizza type : %v and cost is : %v\n", p.Name(), p.Cost())
}

func main() {

	nyPizzaFactory := &nyPizzaFactory{}

	pizza := nyPizzaFactory.CreatePizza()
	display(pizza)

	// chicagoPizzaCreator := chicagoPizzaCreator()
	// pizza1 := chicagoPizzaCreator()
	// display(pizza1)
	factoryMap := make(factoryMap)
	factoryMap["chicago"] = chicagoPizzaCreator

	chicagoCreatorFunc := factoryMap["chicago"]
	pizza1 := chicagoCreatorFunc(3)
	display(pizza1)
}
