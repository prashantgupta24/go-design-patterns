package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/statsd"
	"github.com/pkg/errors"
)

var wg sync.WaitGroup

func foo(c chan int, someValue int) {
	defer wg.Done()
	time.Sleep(time.Second * 5)
	c <- someValue * 5

}

func getVal() string {
	return "abc"
}

func getVal1(index int) string {
	return "abc" + strconv.Itoa(index)
}

type Animal interface {
	Speak() string
}
type Dog struct {
	Name, Sound string
}

func (d Dog) Speak() string {
	return d.Sound
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow"
}

func MakeAnimalSpeak(a Animal) {
	fmt.Println(a.Speak())
}

var test map[int]string
var testInt int
var failedTrainerConnectivityCounter metrics.Counter

func initTest() {
	test = make(map[int]string)
	test[1] = "b"
	testInt := 3
	fmt.Println(testInt)
	s := statsd.New("d", nil)
	failedTrainerConnectivityCounter = s.NewCounter("request_duration_MyMethod_200", 3)
}

func f() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	fmt.Println("Calling g.")
	return g(4)
}

func g(i int) error {
	if i > 3 {
		fmt.Println("Panicking!")
		//panic(fmt.Sprintf("%v", i))
		//log.Fatal("dd")
		//panic(fmt.Errorf("custom error with %v", i))
		// b := 0
		// v := 5 / b
		// fmt.Println(v)
		// err := errors.New("custom error")
		// return errors.Wrap(err, "custom error")
		return fmt.Errorf("custom error with %v", i)
	}
	defer fmt.Println("Defer in g", i)
	fmt.Println("Printing in g", i)
	//g(i + 1)
	return nil
}

func A() error {
	return errors.New("NullPointerException")
}

func B() error {
	return A()
}

type jobMonitorMetrics struct {
	FailedETCDConnectivityCounter        metrics.Counter
	FailedK8sConnectivityCounter         metrics.Counter
	InsufficientK8sResourcesErrorCounter metrics.Counter
	FailedImagePullK8sErrorCounter       metrics.Counter
	FailedETCDWatchCounter               metrics.Counter
	FailedTrainerConnectivityCounter     metrics.Counter
}

type configs struct {
	a      *int
	b      *int
	health func()
}

type options func(*configs)

func New(opts ...options) {
	c := &configs{}
	for _, opt := range opts {
		opt(c)
	}
	c.health()
	fmt.Println(c)
}

func WithHealthFun(fn func()) func(c *configs) {
	return func(c *configs) {
		c.health = fn
	}
}

func BuilderNew() *configs {
	return &configs{}
}

func (c *configs) WithA(aVal int) *configs {
	c.a = &aVal
	return c
}
func (c *configs) WithHealth(fn func()) *configs {
	c.health = fn
	return c
}

func clos() func(string) bool {

	cacheMap := make(map[string]int)
	cacheMap["a"] = 1

	return func(s string) bool {
		_, ok := cacheMap[s]
		if ok {
			return true
		}
		cacheMap[s] = 1
		return false
	}
}

func nonClos(s string) bool {
	cacheMap := make(map[string]int)
	_, ok := cacheMap[s]
	if ok {
		return true
	}
	cacheMap[s] = 1
	return false
}

func sendLoop(sender func() bool) {
	for {
		retry := sender()
		if !retry {
			return
		}
		fmt.Println("inside loop")
		time.Sleep(time.Second)
	}
}

func mySender() func() bool {
	i := 1
	return func() bool {
		fmt.Println("Value of i is :", i)
		if i > 5 {
			return false
		}
		i++
		return true
	}
}

func mySender1() func() bool {
	i := 1
	return func() bool {
		fmt.Println("Value of i is :", i)
		if i > 3 {
			return false
		}
		i++
		return true
	}
}

func addN(m int) func(int) int {
	return func(n int) int {
		return m + n
	}
}

func main() {

	a := addN(5)
	fmt.Println(a(4))
	//*******************************
	// ms := mySender1()
	// sendLoop(ms)

	//*******************************
	// c := clos()
	// fmt.Println(c("aa"))
	// fmt.Println(c("aa"))

	// fmt.Println(nonClos("aa"))
	// fmt.Println(nonClos("aa"))
	//*******************************
	//Option 1
	// fn1 := func(c *configs) {
	// 	fmt.Println("Func 1")
	// 	a := 1
	// 	c.a = &a
	// }

	// fn2 := func() {
	// 	fmt.Println("inside health function")

	// 	// i := 1
	// 	// return func(c *configs) {
	// 	// 	if i > 2 {
	// 	// 		c.health = fn
	// 	// 	} else {
	// 	// 		c.health = func() {
	// 	// 			fmt.Println("Inside custom health function")
	// 	// 		}
	// 	// 	}
	// 	// }
	// }

	// New(fn1, WithHealthFun(fn2))

	//Option 2
	// c := BuilderNew().WithA(2).WithHealth(func() {
	// 	fmt.Println("inside health function")
	// })
	// c.health()
	//fmt.Println(*c.health())

	//*******************************
	//test123()
	// fmt.Println("Main")
	// ABC = "Sdf"

	//*******************************
	// statsdClient := statsd.New(fmt.Sprintf("%s.", "monitor"), log.NewNopLogger())
	// jmMetrics := &jobMonitorMetrics{
	// 	FailedETCDConnectivityCounter:        statsdClient.NewCounter("jobmonitor.etcd.connectivity.failed", 1),
	// 	FailedK8sConnectivityCounter:         statsdClient.NewCounter("jobmonitor.k8s.connectivity.failed", 1),
	// 	InsufficientK8sResourcesErrorCounter: statsdClient.NewCounter("jobmonitor.k8s.insufficientResources.failed", 1),
	// 	FailedImagePullK8sErrorCounter:       statsdClient.NewCounter("jobmonitor.k8s.imagePull.failed", 1),
	// 	FailedETCDWatchCounter:               statsdClient.NewCounter("jobmonitor.etcd.watch.failed", 1),
	// 	FailedTrainerConnectivityCounter:     statsdClient.NewCounter("jobmonitor.trainer.connectivity.failed", 1),
	// }

	// //fmt.Printf("type %v value %v \n\n", reflect.TypeOf(jmMetrics), reflect.ValueOf(jmMetrics).Elem())
	// j := reflect.ValueOf(jmMetrics).Elem()
	// //fmt.Println(reflect.ValueOf(jmMetrics).Elem().Kind())

	// //var ic metrics.Counter = statsdClient.NewCounter("jobmonitor.etcd.connectivity.failed", 1)
	// //fmt.Println(reflect.ValueOf(&ic).Elem().Type())
	// //typeMetricCounter := reflect.ValueOf(&ic).Elem().Type()
	// var typeMetricCounter metrics.Counter
	// fmt.Println(reflect.TypeOf(&typeMetricCounter).Elem())
	// fmt.Println(reflect.TypeOf((*metrics.Counter)(nil)).Elem())

	// for i := 0; i < j.NumField(); i++ {
	// 	// fmt.Println(j.Field(i).Type())
	// 	// fmt.Println(reflect.ValueOf(j.Field(i)).Kind())
	// 	//fmt.Println(j.Field(i).Interface())
	// 	f := j.Field(i)
	// 	fmt.Printf("%d: %s %s = %v\n", i,
	// 		j.Type().Field(i).Name, f.Type(), f.Interface())
	// 	//var i interface {}

	// 	if j.Field(i).Type() == reflect.TypeOf((*metrics.Counter)(nil)).Elem() {
	// 		counter := j.Field(i)
	// 		//fmt.Println(reflect.ValueOf(counter))
	// 		c := counter.Interface().(metrics.Counter)
	// 		c.Add(1)
	// 	}
	// }

	// x := Dog{
	// 	Name:  "Joey",
	// 	Sound: "woof",
	// }
	// v := reflect.ValueOf(&x).Elem()
	// fmt.Println("type:", v.Type())

	// for i := 0; i < v.NumField(); i++ {
	// 	f := v.Field(i)
	// 	if f.Kind() == reflect.String {
	// 		f.SetString("asdf")
	// 	}
	// 	fmt.Println(f)

	// }

	// var x float64 = 3.4
	// v := reflect.ValueOf(&x)
	// v.Elem().SetFloat(7.1)
	// fmt.Println(x)
	//*******************************
	//fmt.Printf("Error: %+v", B())
	//*******************************
	// if err := f(); err != nil {
	// 	log.WithFields(log.Fields{
	// 		"animal": "walrus",
	// 	}).Fatalf("A walrus appears")
	// }
	// fmt.Println("Returned normally from f.")

	//*******************************
	// fooVal := make(chan int)

	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go foo(fooVal, i)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(fooVal)

	// }()
	// for item := range fooVal {
	// 	fmt.Println(item)
	// }
	//*******************************
	// done := make(chan string, 1)

	// go func() {
	// 	//time.Sleep(time.Second * 2)
	// 	done <- getVal()
	// }()

	// if res := <-done; res != "" {
	// 	fmt.Println(res)
	// }

	// fmt.Println("executing here")

	//*******************************

	// var wg sync.WaitGroup

	// for index := 0; index < 10; index++ {
	// 	result := make(chan string)
	// 	i := index
	// 	wg.Add(1)

	// 	go func() {
	// 		result <- getVal1(i)
	// 		close(result)
	// 	}()

	// 	go func() {
	// 		defer wg.Done()
	// 		for val := range result {
	// 			fmt.Println(val)
	// 		}
	// 	}()
	// }

	// wg.Wait()
	// fmt.Println("executing here")
	//*******************************

	// d1 := Dog{
	// 	Name:  "Dog1",
	// 	Sound: "Woof1",
	// }
	// c1 := Cat{
	// 	Name: "Cat1",
	// }
	// MakeAnimalSpeak(d1)
	// MakeAnimalSpeak(c1)

	//*******************************
	// t := time.NewTicker(1 * time.Second)

	// go func() {
	// 	for range t.C {
	// 		fmt.Println("Ticking!")
	// 	}
	// 	fmt.Println("this is never printed")
	// }()

	// time.Sleep(time.Second * 3)
	// fmt.Println("Sending command to close ticker")
	// t.Stop()
	// time.Sleep(time.Second * 3)
	//*******************************

	// fmt.Println(testInt)
	// //initTest()
	// failedTrainerConnectivityCounter.Add(1)
	// fmt.Println(testInt)
	// test[0] = "a"
	// for key, val := range test {
	// 	fmt.Printf("key is %v and val is %v", key, val)
	// 	fmt.Println()
	// }

	//*******************************

}
