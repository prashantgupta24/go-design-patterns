package main

import (
	"fmt"
	"sync"
)

var singletonInstance *singleton

type singleton struct {
	val   int
	mutex sync.RWMutex
}

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
	//fmt.Printf("Got : %v. Now %v : \n", or, s.getVal())
}

func getInstance() *singleton {
	return singletonInstance
}

func main() {
	fmt.Println("Welcome to singleton pattern")
}
