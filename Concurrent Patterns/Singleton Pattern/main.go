package main

import (
	"fmt"
	"sync"
)

var singletonInstance *singleton

type singleton struct {
	val    int
	mutex  sync.RWMutex
	input  chan int
	output chan int
}

//Singleton
func getInstance() *singleton {
	initInstance()
	return singletonInstance
}

//Mutex implementation
func (s *singleton) setVal(val int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.val = val
}

func (s *singleton) getVal() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.val
}

func (s *singleton) addOne() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.val = s.val + 1
}

//Channel implementation
func (s *singleton) addOneThroughChan() {
	s.input <- 1
}

func (s *singleton) getValThroughChan() int {
	s.output <- 1
	val := <-s.output
	return val
}

func initInstance() {
	if singletonInstance == nil {
		singletonInstance = &singleton{
			input:  make(chan int),
			output: make(chan int),
		}

		go func(singletonInstance *singleton) {
			for {
				select {
				case value := <-singletonInstance.input:
					singletonInstance.val = singletonInstance.val + value
				case <-singletonInstance.output:
					singletonInstance.output <- singletonInstance.val
				}
			}
		}(singletonInstance)
	}
}

func main() {
	fmt.Println("Welcome to singleton pattern")
}
