package main

import (
	"fmt"
	"sync"
)

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

var instance singleton

func getInstance() *singleton {
	return &instance
}

func main() {

	var wg sync.WaitGroup

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		fmt.Println("func1")
		defer wg.Done()
		s1 := getInstance()
		s1.setVal(s1.getVal() + 3)
	}(&wg)

	go func(wg *sync.WaitGroup) {
		fmt.Println("func2")
		defer wg.Done()
		s2 := getInstance()
		s2.setVal(s2.getVal() + 2)
	}(&wg)

	wg.Wait()
	s3 := getInstance()
	fmt.Println(s3.getVal())
}
