package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	sleepTime        = 2
	numLoops         = 1
	randomValueLimit = 10
)

type testCase struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testCase))
}

func (suite *testCase) TestFuture1() {
	t := suite.T()

	var successes int
	var failures int

	successCh := make(chan int)
	failureCh := make(chan int)
	resultCh := make(chan int)

	go func() {
		for {
			select {
			case <-successCh:
				successes++
			case <-failureCh:
				failures++
			case <-resultCh:
				resultCh <- successes + failures
				return
			}
		}
	}()

	for i := 0; i < numLoops; i++ {
		val := rand.Intn(randomValueLimit)

		NewFuture().Success(func(str string) {
			fmt.Println("func1 was successful")
			successCh <- 1
		}).Fail(func(err error) {
			fmt.Printf("error returned from func1 %v\n", err)
			failureCh <- 1
		}).Execute(funcToExecute1, val)

		NewFuture().Success(func(str string) {
			fmt.Println("func2 was successful")
			successCh <- 1
		}).Fail(func(err error) {
			fmt.Printf("error returned from func2 %v\n", err)
			failureCh <- 1
		}).Execute(funcToExecute2, val)

		func3Wrap := func(int) (string, error) {
			return funcToExecute3()
		}
		NewFuture().Success(func(str string) {
			fmt.Println("func3 was successful")
			successCh <- 1
		}).Fail(func(err error) {
			fmt.Printf("error returned from func3 %v\n", err)
			failureCh <- 1
		}).Execute(func3Wrap, val)
	}

	time.Sleep(time.Second * time.Duration(sleepTime+1))

	resultCh <- 1 //activate result channel to send us the result
	totalOutcomes := <-resultCh
	assert.Equal(t, numLoops*3, totalOutcomes, "success and failures should be %v, instead it is %v", numLoops*3, totalOutcomes)
}

func sleep() {
	time.Sleep(time.Second * time.Duration(rand.Intn(sleepTime)))
}

func funcToExecute1(val int) (string, error) {
	//fmt.Println("func1")
	sleep()
	if val > 10 {
		return "", fmt.Errorf("value too high for func1! : %v", val)
	}
	return "value is correct for func1!", nil
}

func funcToExecute2(val int) (string, error) {
	//fmt.Println("func2")
	sleep()
	if val%2 != 0 {
		return "", fmt.Errorf("value not divisible by 2! : %v", val)
	}
	return "value is correct for func2!", nil
}

func funcToExecute3() (string, error) {
	//fmt.Println("func3")
	sleep()
	return "func3 always passes!", nil
}
