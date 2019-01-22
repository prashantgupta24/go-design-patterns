package main

import "fmt"

type Subject interface {
	Subscribe(Observer)
	Unsubscribe(Observer)
	OnChange()
}

type Observer interface {
	Update(int)
}

type WeatherStation struct {
	observers []Observer
	data      int
}

func (w *WeatherStation) Subscribe(o Observer) {
	w.observers = append(w.observers, o)
}

func (w *WeatherStation) Unsubscribe(o Observer) {

}

func (w *WeatherStation) OnChange() {
	for _, observer := range w.observers {
		observer.Update(w.data)
	}
}

func (w *WeatherStation) SetVal(value int) {
	w.data = value
	w.OnChange()
}

type DisplayOne struct{}
type DisplayTwo struct{}
type DisplayThree struct{}

func (d *DisplayOne) Update(val int) {
	fmt.Println("Updating display one by : ", val)
}
func (d *DisplayTwo) Update(val int) {
	fmt.Println("Updating display two by : ", val)
}
func (d *DisplayThree) Update(val int) {
	fmt.Println("Updating display three by : ", val)
}

func main() {
	fmt.Println("hello")

	weatherStation := &WeatherStation{}

	d1 := &DisplayOne{}
	d2 := &DisplayTwo{}
	d3 := &DisplayThree{}

	weatherStation.Subscribe(d1)
	weatherStation.Subscribe(d2)
	weatherStation.Subscribe(d3)

	weatherStation.SetVal(2)
}
