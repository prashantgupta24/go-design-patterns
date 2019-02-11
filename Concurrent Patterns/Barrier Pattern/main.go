package main

import (
	"sync"
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

//add all jobs to a new instance of barrier
barrier := NewBarrier().AddN("job1", job1).
	AddN("job2", job2).AddN("job3", job3)

Option 1:
results, err := barrier.Execute()

Execute() returns a Go or no-go, i.e. if there was an error
in any of the jobs submitted, that error is returned.
If all jobs passed, then all the results are returned
as a map of the job name and their corresponding result.

We can just fetch the result of a function by querying the
response map returned:

//Result of Job 1 (assuming all jobs passed)
job1Output := results["job1"]

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
	resultsCh chan *Result
	functions map[string]functionType
}

//NewBarrier creates a new Barrier
func NewBarrier() *Barrier {
	b := &Barrier{}
	b.init()
	return b
}

//Result is the information being passed back to the user.
type Result struct {
	funcName     string
	funcResponse func() interface{}
	err          error
}

//initializes the Barrier struct, called automatically
func (b *Barrier) init() {
	var wg sync.WaitGroup
	b.wg = &wg
	b.resultsCh = make(chan *Result)
	b.functions = make(map[string]functionType)
}

//all functions need to be of this type
type functionType func(int) (func() interface{}, error)

//Add adds a function to our Barrier execution queue
func (b *Barrier) Add(fn functionType) *Barrier {
	//b.functions = append(b.functions, fn)
	return b.AddN("default", fn)
}

/*AddN adds a function to our Barrier execution queue,
along with a name to the function. This can be used to fetch
the corresponding result of the function
*/
func (b *Barrier) AddN(functionName string, fn functionType) *Barrier {
	b.functions[functionName] = fn
	return b
}

/*ExecuteAndReturnResults returns an array of results for the user
to handle. Needed if the returned errors need to be handled
separately

Also see execute()
*/
func (b *Barrier) ExecuteAndReturnResults(val int) map[string]*Result {
	b.executeDefault(&val)
	return b.wait()
}

/*Execute parses the array of results, and returns all results if no error,
else an error if any one of the jobs failed.

Also see executeAndReturn()
*/
func (b *Barrier) Execute(val int) (map[string]*Result, error) {
	b.executeDefault(&val)
	results := b.wait()

	for _, result := range results {
		if result.err != nil {
			return nil, result.err
		}
	}
	return results, nil
}

//executeDefault is not a public function
func (b *Barrier) executeDefault(val *int) {
	for name, fn := range b.functions {
		b.wg.Add(1)
		go func(fn functionType, name string, b *Barrier, val *int) {
			defer b.wg.Done()
			resp, err := fn(*val)
			if err != nil {
				b.resultsCh <- &Result{
					funcName:     name,
					funcResponse: nil,
					err:          err,
				}
			} else {
				b.resultsCh <- &Result{
					funcName:     name,
					funcResponse: resp,
					err:          nil,
				}
			}
		}(fn, name, b, val)
	}
}

//wait is not a public function
func (b *Barrier) wait() map[string]*Result {
	go func(wg *sync.WaitGroup, results chan *Result) {
		wg.Wait()
		close(results)
	}(b.wg, b.resultsCh)

	results := make(map[string]*Result)
	for result := range b.resultsCh {
		results[result.funcName] = result
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
