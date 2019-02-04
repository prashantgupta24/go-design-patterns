package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testCase struct {
	suite.Suite
	barrier *barrier
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testCase))
}

//Run once before all tests
func (suite *testCase) SetupSuite() {
	barrier := &barrier{}

	//since job3 doesn't require an input, we have to wrap it around a function
	job3Wrapper := func(int) (string, error) {
		return job3()
	}

	barrier.add(job1).add(job2).add(job3Wrapper)
	suite.barrier = barrier
}

func job1(val int) (string, error) {
	fmt.Println("executing job1")
	time.Sleep(time.Second * 3)
	if val > 10 {
		return "success", nil
	}

	errMsg := fmt.Sprintf("too less for val in func1 : %v. It needs greater than 10 ", val)
	return "", customErrorNew(errMsg, false)
}

func job2(val int) (string, error) {
	fmt.Println("executing job2")
	time.Sleep(time.Second * 2)
	if val%2 == 0 {
		return "success", nil
	}
	errMsg := fmt.Sprintf("CRITICAL ERROR!! Val not divisible by 2 in func2 : %v", val)
	return "", customErrorNew(errMsg, true)
}

func job3() (string, error) {
	fmt.Println("executing job3")
	time.Sleep(time.Second * 2)
	return "func3 always passes!", nil
}
func (suite *testCase) TestExecute1() {

	barrier := suite.barrier
	t := suite.T()

	_, err := barrier.execute(11)

	assert.NotNil(t, err, "Error should not be nil. ")
	assert.IsType(t, err, &customError{}, "Should be of custom error type")

	customErr := err.(*customError)
	assert.True(t, customErr.critical, "It should be critical")
}

func (suite *testCase) TestExecute2() {

	barrier := suite.barrier
	t := suite.T()

	_, err := barrier.execute(2)

	assert.NotNil(t, err, "Error should not be nil. ")
	assert.IsType(t, err, &customError{}, "Should be of custom error type")

	customErr := err.(*customError)
	assert.False(t, customErr.critical, "It should not be critical")
}

func (suite *testCase) TestExecute3() {

	barrier := suite.barrier
	t := suite.T()

	_, err := barrier.execute(12)

	assert.Nil(t, err, "Error should be nil. ", err)
}
