package main

import (
	"fmt"
)

//all functions need to be of this type
type functionType func(int) (string, error)

type successFunction func(string)
type errorFunction func(error)

type future struct {
	successFunc successFunction
	failFunc    errorFunction
}

//NewFuture creates a new future
func NewFuture() *future {
	return &future{}
}

func (f *future) execute(fn functionType, val int) *future {
	f.init()
	go func(f *future) {
		if str, err := fn(val); err != nil {
			f.failFunc(err)
		} else {
			f.successFunc(str)
		}
	}(f)
	return f
}

func (f *future) init() {
	if f.failFunc == nil {
		f.failFunc = func(err error) {
			fmt.Println("ERROR !!! >> ", err)
		}
	}
	if f.successFunc == nil {
		f.successFunc = func(str string) {
			fmt.Println(str)
		}
	}
}

func (f *future) success(fn successFunction) *future {
	f.successFunc = fn
	return f
}

func (f *future) fail(fn errorFunction) *future {
	f.failFunc = fn
	return f
}

func main() {

}
