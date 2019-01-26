package main

import (
	"fmt"
	"sync"
	"time"
)

type str struct {
	val1 int
	val2 int
}

type returnVal struct {
	msg string
	err error
}

type barrier struct {
	wg        *sync.WaitGroup
	results   chan *returnVal
	functions []func()
}

func (b *barrier) init() {
	var wg sync.WaitGroup
	b.wg = &wg
	b.results = make(chan *returnVal)
}

func (b *barrier) add(fn func()) {
	b.functions = append(b.functions, fn)
}

func (b *barrier) execute() (string, error) {

	for _, fn := range b.functions {
		b.wg.Add(1)
		go fn()
	}
	// wg.Add(1)
	// go job1(val1, &wg, results, st)

	// wg.Add(1)
	// go job2(val2, &wg, results, st)

	returnValues := b.wait()

	for _, returnValue := range returnValues {
		if returnValue.err != nil {
			return "", returnValue.err
		}
	}
	return fmt.Sprintf("Values are correct!"), nil
}

// func (b *barrier) wait() (string, error) {
// 	go func(wg *sync.WaitGroup, results chan *returnVal) {
// 		wg.Wait()
// 		close(results)
// 	}(b.wg, b.results)

// 	for result := range b.results {
// 		if result.err != nil {
// 			return "", result.err
// 		}
// 	}
// 	return fmt.Sprintf("Correct!"), nil
// }

func (b *barrier) wait() []*returnVal {
	go func(wg *sync.WaitGroup, results chan *returnVal) {
		wg.Wait()
		close(results)
	}(b.wg, b.results)

	var returnValues []*returnVal
	for result := range b.results {
		returnValues = append(returnValues, result)
	}
	return returnValues
}
func job1(barrier *barrier, val int, st *str) {
	defer barrier.wg.Done()

	time.Sleep(time.Second * 3)
	if val > 1 {
		st.val1 = val
		barrier.results <- &returnVal{
			msg: "success",
			err: nil,
		}
	} else {
		err := fmt.Errorf("too much for val in func1 : %v", val)
		barrier.results <- &returnVal{
			msg: "",
			err: err,
		}
	}
}

func job2(barrier *barrier, val int, st *str) {
	defer barrier.wg.Done()

	time.Sleep(time.Second * 2)
	if val < 0 {
		st.val2 = val
		barrier.results <- &returnVal{
			msg: "success",
			err: nil,
		}
	} else {
		err := fmt.Errorf("too less for val in func2: %v", val)
		barrier.results <- &returnVal{
			msg: "",
			err: err,
		}
	}
}

func job3() (string, error) {
	time.Sleep(time.Second * 2)
	return "func3 always passes!", nil
}

func createJobs(val1, val2 int) (string, error) {
	st := &str{}

	barrier := &barrier{}
	barrier.init()

	barrier.wg.Add(1)
	go job1(barrier, val1, st)

	barrier.wg.Add(1)
	go job2(barrier, val2, st)

	barrier.wg.Add(1)
	go func() {
		defer barrier.wg.Done()
		str, err := job3()
		if err != nil {
			err := fmt.Errorf("incorrect")
			barrier.results <- &returnVal{
				msg: "",
				err: err,
			}
		} else {
			barrier.results <- &returnVal{
				msg: str,
				err: nil,
			}
		}

	}()

	returnValues := barrier.wait()

	for _, returnValue := range returnValues {
		if returnValue.err != nil {
			return "", returnValue.err
		}
	}
	return fmt.Sprintf("Values are correct! Struct is: %v", st), nil
	// var wg sync.WaitGroup
	// results := make(chan *returnVal, 2)

	// wg.Add(1)
	// go job1(val1, &wg, results, st)

	// wg.Add(1)
	// go job2(val2, &wg, results, st)

	// go func(wg *sync.WaitGroup, results chan *returnVal) {
	// 	wg.Wait()
	// 	close(results)
	// }(&wg, results)

	// for result := range results {
	// 	if result.err != nil {
	// 		return "", result.err
	// 	}
	// }
	// return fmt.Sprintf("Values are correct! Struct is: %v", st), nil
}

func display(s string, err error) {
	if err != nil {
		fmt.Println("ERROR!! >>> ", err)
	} else {
		fmt.Println(s)
	}
}

func main() {

	display(createJobs(4, 2))
	display(createJobs(-4, -2))
	display(createJobs(-4, 2))
	display(createJobs(4, -2))

}
