package main

import (
	"fmt"
)

/*

The Future design pattern (also called Promise) is a quick and easy
way to achieve concurrent structures for asynchronous programming

The idea here is to achieve a fire-and-forget that handles all possible
results in an action. In short, we will define each possible behavior
of an action before executing them in different Goroutines.

Basically define the behavior in advance and let the future resolve
the possible solutions.

Usage:

NewFuture().Success(func(str string) {
			...
		}).Fail(func(err error) {
			...
		}).Execute(funcToExecute)

We just fire and forget. After the function executes, it is the
responsibilty of the Future to make sure to run the success function
if the execution was successful, if not run the fail
function. We do not wait for the result.

Note: This is different from the barrier pattern in which
we also fire concurrent jobs but then wait for all of them to
finish at the end.

*/

//Future struct is the main struct that needs to be initialized
type Future struct {
	successFunc successFunction
	failFunc    errorFunction
}

//NewFuture creates a new future
func NewFuture() *Future {
	return &Future{}
}

//all functions need to be of this type
type functionType func(int) (string, error)

type successFunction func(string)
type errorFunction func(error)

/*
Execute executes the main function, and correspondingly
handles the result
*/
func (f *Future) Execute(fn functionType, val int) *Future {
	f.init()
	go func(f *Future) {
		if str, err := fn(val); err != nil {
			f.failFunc(err)
		} else {
			f.successFunc(str)
		}
	}(f)
	return f
}

func (f *Future) init() {
	if f.failFunc == nil {
		f.failFunc = func(err error) {
			panic(err)
		}
	}
	if f.successFunc == nil {
		f.successFunc = func(str string) {
			fmt.Println(str)
		}
	}
}

//Success function to run if the execution was successful
func (f *Future) Success(fn successFunction) *Future {
	f.successFunc = fn
	return f
}

//Fail function to run if the execution was successful
func (f *Future) Fail(fn errorFunction) *Future {
	f.failFunc = fn
	return f
}

func main() {

}
