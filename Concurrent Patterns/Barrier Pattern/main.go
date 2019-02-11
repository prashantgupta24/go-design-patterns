package main

import (
	"fmt"
	"sync"
	"time"
)

/*

Barrier pattern

Its purpose is simple--put up a barrier so that nobody passes
until we have all the results we need, something quite common in
concurrent applications.

Imagine the situation where we have a micro-services application
where one service needs to compose its response by merging the responses
of other microservices. This is where the Barrier pattern can help us.

Our Barrier pattern could be a service that will block its response
until it has been composed with the results returned by one or more
different Goroutines (or services).

Usage:

barrier := &Barrier{}

//add all jobs to barrier
barrier.Add(job1).Add(job2).Add(job3)

Option 1:
resp, err := barrier.Execute()

Execute() simply returns a Go or no-go, i.e. if there was an error
in any of the jobs submitted, that error is returned.
If all jobs passed, then error will be nil.

Option 2:
results := Barrier.executeAndReturnResults()

If we want more control on each of the job's result, then
we can use ExecuteAndReturnResults(), which returns an array
of results for us to deal with.

for _, result := range results {
	...
}

*/

//Barrier struct is the main struct containing all components we need
type Barrier struct {
	wg        *sync.WaitGroup
	results   chan *Result
	functions []functionType
}

//Result is the information being passed back to the user.
type Result struct {
	response func() interface{}
	err      error
}

//initializes the Barrier struct, called automatically
func (b *Barrier) init() {
	var wg sync.WaitGroup
	b.wg = &wg
	b.results = make(chan *Result)
}

//all functions need to be of this type
type functionType func(int) (func() interface{}, error)

//Add adds a function to our Barrier execution queue
func (b *Barrier) Add(fn functionType) *Barrier {
	b.functions = append(b.functions, fn)
	return b
}

/*ExecuteAndReturnResults returns an array of results for the user
to handle.

Also see execute()
*/
func (b *Barrier) ExecuteAndReturnResults(val int) []*Result {
	b.executeDefault(&val)
	return b.wait()
}

/*Execute parses the array of results, and only returns an error
if any one of the jobs failed.

Also see executeAndReturn()
*/
func (b *Barrier) Execute(val int) (string, error) {
	b.executeDefault(&val)
	results := b.wait()

	for _, result := range results {
		if result.err != nil {
			return "", result.err
		}
	}
	return fmt.Sprintf("Values are correct!"), nil
}

//executeDefault is not a public function
func (b *Barrier) executeDefault(val *int) {
	b.init()
	for _, fn := range b.functions {
		b.wg.Add(1)
		go func(fn functionType, b *Barrier, val *int) {
			defer b.wg.Done()
			resp, err := fn(*val)
			if err != nil {
				b.results <- &Result{
					response: nil,
					err:      err,
				}
			} else {
				b.results <- &Result{
					response: resp,
					err:      nil,
				}
			}
		}(fn, b, val)
	}
}

//wait is not a public function
func (b *Barrier) wait() []*Result {
	go func(wg *sync.WaitGroup, results chan *Result) {
		wg.Wait()
		close(results)
	}(b.wg, b.results)

	var results []*Result
	for result := range b.results {
		results = append(results, result)
	}
	return results
}

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

func job12(val int) (func() interface{}, error) {
	fmt.Println("executing job12")
	time.Sleep(time.Second * 3)
	localVal := 10
	if val > 10 {
		return func() interface{} {
			return localVal
		}, nil
	}

	errMsg := fmt.Sprintf("too less for val in func1 : %v. It needs greater than 10 ", val)
	return nil, customErrorNew(errMsg, false)
}

func main() {
	//option1, we only care if any critical errors occured in any of the jobs
	// Barrier := &Barrier{}
	// Barrier.Add(job)
	// resp, err := Barrier.Execute(21)
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
	// Barrier := &Barrier{}
	// Barrier.Add(job12)
	// results := Barrier.ExecuteAndReturnResults(12)

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
	// 	} else {
	// 		resp := result.response()
	// 		switch valType := resp.(type) {
	// 		case int:
	// 			fmt.Println("function returned int : ", valType)
	// 		case bool:
	// 			fmt.Println("function returned bool : ", valType)
	// 		}
	// 	}
	// }
	// if !hasError {
	// 	fmt.Println("Values are correct!")
	// }
}

//TODO custom error, interface input
