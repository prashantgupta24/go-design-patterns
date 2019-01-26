package main

import (
	"fmt"
	"sync"
	"time"
)

type returnVal struct {
	msg string
	err error
}

//all functions need to be of this type
type functionType func(int) (string, error)

type barrier struct {
	wg        *sync.WaitGroup
	results   chan *returnVal
	functions []functionType
}

func (b *barrier) init() {
	var wg sync.WaitGroup
	b.wg = &wg
	b.results = make(chan *returnVal)
}

func (b *barrier) add(fn functionType) *barrier {
	b.functions = append(b.functions, fn)
	return b
}

func (b *barrier) executeAndReturn(val int) []*returnVal {
	b.executeDefault(&val)
	return b.wait()
}

func (b *barrier) executeDefault(val *int) {
	for _, fn := range b.functions {
		b.wg.Add(1)
		go func(fn functionType, b *barrier, val *int) {
			defer b.wg.Done()
			resp, err := fn(*val)
			if err != nil {
				b.results <- &returnVal{
					msg: "",
					err: err,
				}
			} else {
				b.results <- &returnVal{
					msg: resp,
					err: nil,
				}
			}
		}(fn, b, val)
	}
}

func (b *barrier) execute(val int) (string, error) {
	b.executeDefault(&val)
	returnValues := b.wait()

	for _, returnValue := range returnValues {
		if returnValue.err != nil {
			return "", returnValue.err
		}
	}
	return fmt.Sprintf("Values are correct!"), nil
}

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

func job1(val int) (string, error) {
	fmt.Println("executing job1")
	time.Sleep(time.Second * 3)
	if val > 10 {
		return "success", nil
	}
	err := fmt.Errorf("too less for val in func1 : %v. It needs greater than 10 ", val)
	return "", err
}

func job2(val int) (string, error) {
	fmt.Println("executing job2")
	time.Sleep(time.Second * 2)
	if val%2 == 0 {
		return "success", nil
	}
	err := fmt.Errorf("Val not divisible by 2 in func2 : %v", val)
	return "", err
}

func job3(val int) (string, error) {
	fmt.Println("executing job3")
	time.Sleep(time.Second * 2)
	return "func3 always passes!", nil
}

func createJobs(val int) (string, error) {

	barrier := &barrier{}
	barrier.init()

	barrier.add(job1).add(job2).add(job3)

	//option1, only know if error occured in any of the jobs
	//return barrier.execute(val)

	//option2, more control on return values
	returnValues := barrier.executeAndReturn(val)
	for _, returnValue := range returnValues {
		if returnValue.err != nil {
			return "", returnValue.err
		}
	}
	return fmt.Sprintf("Values are correct!"), nil
}

func display(s string, err error) {
	if err != nil {
		fmt.Println("ERROR!! >>> ", err)
	} else {
		fmt.Println(s)
	}
}

func main() {
	display(createJobs(4))
	display(createJobs(12))
}
