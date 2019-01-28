package main

import (
	"fmt"
	"sync"
	"time"
)

//all functions need to be of this type
type functionType func(int) (string, error)

type successFunction func(string)
type errorFunction func(error)

type future struct {
	successFunc successFunction
	failFunc    errorFunction
	wg          *sync.WaitGroup
}

func (f *future) execute(fn functionType, val int) {
	f.init()
	f.wg.Add(1)
	go func(f *future) {
		if str, err := fn(val); err != nil {
			f.failFunc(err)
			f.wg.Done()
		} else {
			f.successFunc(str)
			f.wg.Done()
		}
	}(f)
}

func (f *future) init() {
	if f.wg == nil {
		var wg sync.WaitGroup
		f.wg = &wg
	}
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

func funcToExecute1(val int) (string, error) {
	time.Sleep(time.Second * 2)
	if val > 10 {
		return "", fmt.Errorf("value too high for func1! : %v", val)
	}
	return "value is correct for func1!", nil
}

func funcToExecute2(val int) (string, error) {
	time.Sleep(time.Second * 3)
	if val%2 != 0 {
		return "", fmt.Errorf("value not divisible by 2! : %v", val)
	}
	return "value is correct for func2!", nil
}

func funcToExecute3() (string, error) {
	return "func3 always passes!", nil
}

func main() {

	futures := &future{}

	val := 40

	futures.success(func(str string) {
		fmt.Println(str)
	}).fail(func(err error) {
		fmt.Println("ERROR !!! >> ", err)
	}).execute(funcToExecute1, val)

	futures.success(func(str string) {
		fmt.Println(str)
	}).fail(func(err error) {
		fmt.Println("ERROR !!! >> ", err)
	}).execute(funcToExecute2, val)

	futures.success(func(str string) {
		fmt.Println(str)
	}).fail(func(err error) {
		fmt.Println("ERROR !!! >> ", err)
	}).execute(func(int) (string, error) {
		return funcToExecute3()
	}, 0)

	futures.wg.Wait()
}
