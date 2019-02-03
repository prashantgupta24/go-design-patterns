package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testCase struct {
	suite.Suite
	wg      *sync.WaitGroup
	loopVal int
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testCase))
}

func (suite *testCase) SetupSuite() {
	var wg sync.WaitGroup
	suite.wg = &wg
	suite.loopVal = 1000
}

//Run before each test
func (suite *testCase) BeforeTest(suiteName, testName string) {
	fmt.Println("running before")
	singletonInstance = nil
}

func (suite *testCase) TestSingletonMutex1() {

	t := suite.T()
	wg := suite.wg
	loopVal := suite.loopVal
	s := getInstance()

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < loopVal; i++ {
			s := getInstance()
			s.addOne()
		}
	}(wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < loopVal; i++ {
			s := getInstance()
			s.addOne()
		}
	}(wg)

	wg.Wait()
	assert.Equal(t, loopVal*2, s.getVal(), "Values not equal. %v != %v", s.getVal(), loopVal*2)
}

func (suite *testCase) TestSingletonMutex2() {

	t := suite.T()
	wg := suite.wg
	loopVal := suite.loopVal
	s := getInstance()

	wg.Add(loopVal)
	for i := 0; i < loopVal; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			s := getInstance()
			s.addOne()
		}(wg)
	}

	wg.Wait()
	assert.Equal(t, loopVal, s.getVal(), "Values not equal. %v != %v", s.getVal(), loopVal)
}

func (suite *testCase) TestSingletonChan1() {

	t := suite.T()
	wg := suite.wg
	loopVal := suite.loopVal
	s := getInstance()

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < loopVal; i++ {
			s := getInstance()
			s.addOneThroughChan()
		}
	}(wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < loopVal; i++ {
			s := getInstance()
			s.addOneThroughChan()
		}
	}(wg)

	wg.Wait()

	actualVal := s.getValThroughChan()
	assert.Equal(t, loopVal*2, actualVal, "Values not equal. %v != %v", loopVal*2, actualVal)
}

func (suite *testCase) TestSingletonChan2() {

	t := suite.T()
	wg := suite.wg
	loopVal := suite.loopVal
	s := getInstance()

	wg.Add(loopVal)
	for i := 0; i < loopVal; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			s := getInstance()
			s.addOneThroughChan()
		}(wg)
	}

	wg.Wait()

	actualVal := s.getValThroughChan()
	assert.Equal(t, loopVal, actualVal, "Values not equal. %v != %v", loopVal*2, actualVal)
}
