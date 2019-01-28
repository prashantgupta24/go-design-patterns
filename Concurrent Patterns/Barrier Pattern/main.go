package main

import (
	"fmt"
	"sync"
	"time"
)

type customError struct {
	err      string
	critical bool
}

func (e *customError) Error() string {
	return e.err
}

func customErrorNew(text string, critical bool) error {
	return &customError{
		err:      text,
		critical: critical,
	}
}

type result struct {
	msg string
	err error
}

//all functions need to be of this type
type functionType func(int) (string, error)

type barrier struct {
	wg        *sync.WaitGroup
	results   chan *result
	functions []functionType
}

func (b *barrier) init() {
	var wg sync.WaitGroup
	b.wg = &wg
	b.results = make(chan *result)
}

func (b *barrier) add(fn functionType) *barrier {
	b.functions = append(b.functions, fn)
	return b
}

func (b *barrier) executeAndReturn(val int) []*result {
	b.executeDefault(&val)
	return b.wait()
}

func (b *barrier) executeDefault(val *int) {
	b.init()
	for _, fn := range b.functions {
		b.wg.Add(1)
		go func(fn functionType, b *barrier, val *int) {
			defer b.wg.Done()
			resp, err := fn(*val)
			if err != nil {
				b.results <- &result{
					msg: "",
					err: err,
				}
			} else {
				b.results <- &result{
					msg: resp,
					err: nil,
				}
			}
		}(fn, b, val)
	}
}

func (b *barrier) execute(val int) (string, error) {
	b.executeDefault(&val)
	results := b.wait()

	for _, result := range results {
		if result.err != nil {
			return "", result.err
		}
	}
	return fmt.Sprintf("Values are correct!"), nil
}

func (b *barrier) wait() []*result {
	go func(wg *sync.WaitGroup, results chan *result) {
		wg.Wait()
		close(results)
	}(b.wg, b.results)

	var results []*result
	for result := range b.results {
		results = append(results, result)
	}
	return results
}

func job1(val int) (string, error) {
	fmt.Println("executing job1")
	time.Sleep(time.Second * 3)
	if val > 10 {
		return "success", nil
	}

	errMsg := fmt.Sprintf("too less for val in func1 : %v. It needs greater than 10 ", val)
	return "", customErrorNew(errMsg, false)
}

func job2(val int) (string, error) {
	fmt.Println("executing job2")
	time.Sleep(time.Second * 2)
	if val%2 == 0 {
		return "success", nil
	}
	errMsg := fmt.Sprintf("Val not divisible by 2 in func2 : %v", val)
	return "", customErrorNew(errMsg, true)
}

func job3() (string, error) {
	fmt.Println("executing job3")
	time.Sleep(time.Second * 2)
	return "func3 always passes!", nil
}

func handleJobs(val int) {

	barrier := &barrier{}

	job3Wrapper := func(int) (string, error) {
		return job3()
	}

	barrier.add(job1).add(job2).add(job3Wrapper)

	//option1, we only care if any critical errors occured in any of the jobs
	resp, err := barrier.execute(val)
	if err != nil {
		if err, ok := err.(*customError); ok {
			if err.critical {
				fmt.Println("CRITICAL ERROR!! ", err)
			} else {
				fmt.Println("ERROR!! ", err)
			}
		} else {
			fmt.Println("ERROR!! >>> ", err)
		}
	} else {
		//SUCCESS, all jobs passed without errors
		fmt.Println(resp)
	}

	//option2, if we need more control, we get the list of results
	// and execute based on each result
	// results := barrier.executeAndReturn(val)

	// hasError := false
	// for _, result := range results {
	// 	if result.err != nil {
	// 		hasError = true
	// 		if err, ok := result.err.(*customError); ok {
	// 			if err.critical {
	// 				fmt.Println("CRITICAL ERROR!! ", err)
	// 			} else {
	// 				fmt.Println("ERROR!! ", err)
	// 			}
	// 		} else {
	// 			fmt.Println("ERROR!! >>> ", result.err)
	// 		}
	// 	}
	// }
	// if !hasError {
	// 	fmt.Println("Values are correct!")
	// }
}

func main() {
	handleJobs(4)
	handleJobs(11)
	handleJobs(12)
}

//TODO custom error, interface input
