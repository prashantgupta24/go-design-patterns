package main

import (
	"fmt"
	"sync"
)

/*

Barrier pattern

Its purpose is simple--put up a barrier so that nobody passes
until we have all the results we need, something quite common in concurrent applications.

Imagine the situation where we have a microservices application
where one service needs to compose its response by merging the responses
of other microservices. This is where the Barrier pattern can help us.

Our Barrier pattern could be a service that will block its response
until it has been composed with the results returned by one or more
different Goroutines (or services).

Usage:

barrier := &barrier{}

//add all jobs to barrier
barrier.add(job1).add(job2).add(job3Wrapper)

resp, err := barrier.execute()
//handle the error as you see fit

*/

//CUSTOM ERROR SECTION

//custom error can be returned from the
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

//each result looks like this
type result struct {
	msg string
	err error
}

//all functions need to be of this type
type functionType func(int) (string, error)

//main struct
type barrier struct {
	wg        *sync.WaitGroup
	results   chan *result
	functions []functionType
}

//initializes the barrier struct
func (b *barrier) init() {
	var wg sync.WaitGroup
	b.wg = &wg
	b.results = make(chan *result)
}

//adds a function to our barrier struct
func (b *barrier) add(fn functionType) *barrier {
	b.functions = append(b.functions, fn)
	return b
}

//executeAndReturn returns an array of results for the user to handle. Also see execute()
func (b *barrier) executeAndReturn(val int) []*result {
	b.executeDefault(&val)
	return b.wait()
}

//executeDefault is not a public function
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

/*execute parses the array of results, and only returns an error
if any one of the jobs failed. Also see executeAndReturn()
*/
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

//wait is not a public function
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

func main() {
	//option1, we only care if any critical errors occured in any of the jobs
	//resp, err := barrier.execute(val)
	// if err != nil {
	// 	if err, ok := err.(*customError); ok {
	// 		if err.critical {
	// 			fmt.Println("CRITICAL ERROR!! ", err)
	// 		} else {
	// 			fmt.Println("ERROR!! ", err)
	// 		}
	// 	} else {
	// 		fmt.Println("ERROR!! >>> ", err)
	// 	}
	// } else {
	// 	//SUCCESS, all jobs passed without errors
	// 	fmt.Println(resp)
	// }

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

//TODO custom error, interface input
