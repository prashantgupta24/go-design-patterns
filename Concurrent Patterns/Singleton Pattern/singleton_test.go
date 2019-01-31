package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testCase struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testCase))
}

//Run once before each test
func (suite *testCase) SetupTest() {

	//create a fresh singleton instance for each test
	singletonInstance = &singleton{}
}

func (suite *testCase) TestSingleton1() {

	t := suite.T()

	var wg sync.WaitGroup
	loopVal := 1000

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < loopVal; i++ {
			s := getInstance()
			s.addOne()
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < loopVal; i++ {
			s := getInstance()
			s.addOne()
		}
	}(&wg)

	wg.Wait()
	assert.Equal(t, getInstance().getVal(), loopVal*2, "Values not equal. %v != %v", getInstance().getVal(), loopVal*2)
}

func (suite *testCase) TestSingleton2() {

	t := suite.T()

	var wg sync.WaitGroup
	loopVal := 1000

	wg.Add(loopVal)
	for i := 0; i < loopVal; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			s := getInstance()
			s.addOne()
		}(&wg)
	}

	wg.Wait()
	assert.Equal(t, getInstance().getVal(), loopVal, "Values not equal. %v != %v", getInstance().getVal(), loopVal)
}
