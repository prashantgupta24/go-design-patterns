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
	barrier *Barrier
}

type testCaseCustom struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testCase))
	suite.Run(t, new(testCaseCustom))
}

//Run once before all tests
func (suite *testCase) SetupSuite() {
	//since job3 doesn't require an input, we have to wrap it around a function
	job3Wrapper := func(int) (func() interface{}, error) {
		return job3()
	}

	barrier := NewBarrier().Add(job1).Add(job2).Add(job3Wrapper)
	suite.barrier = barrier
}

func (suite *testCase) TestExecute1() {

	barrier := suite.barrier
	t := suite.T()

	_, err := barrier.Execute(11)

	assert.NotNil(t, err, "Error should not be nil. ")
	assert.IsType(t, err, &customError{}, "Should be of custom error type")

	customErr := err.(*customError)
	assert.True(t, customErr.critical, "It should be critical")
}

func (suite *testCase) TestExecute2() {

	barrier := suite.barrier
	t := suite.T()

	_, err := barrier.Execute(2)

	assert.NotNil(t, err, "Error should not be nil. ")
	assert.IsType(t, err, &customError{}, "Should be of custom error type")

	customErr := err.(*customError)
	assert.False(t, customErr.critical, "It should not be critical")
}

func (suite *testCase) TestExecute3() {

	barrier := suite.barrier
	t := suite.T()

	_, err := barrier.Execute(12)

	assert.Nil(t, err, "Error should be nil. ", err)
}

func (suite *testCaseCustom) TestCustom() {
	t := suite.T()

	barrier := NewBarrier().AddN("customJob1", customJob1).
		AddN("customJob2", customJob2).AddN("customJob3", customJob3)

	input := 12
	results, err := barrier.Execute(input)
	assert.Nil(t, err, "err should be nil")

	if err != nil {
		if err, ok := err.(*customError); ok {
			if err.critical {
				fmt.Println("CRITICAL ERROR!! ", err)
			} else {
				fmt.Println("ERROR!! ", err)
			}
		} else {
			fmt.Println("ERROR!! >>> ", err)
		}
	} else {
		resultForCustomJob1 := results["customJob1"]
		resp1 := resultForCustomJob1.funcResponse()
		fmt.Println("customJob1 returned : ", resp1)
		assert.Equal(t, input*10, resp1, "custom job 1 did not match")

		resultForCustomJob2 := results["customJob2"]
		resp2 := resultForCustomJob2.funcResponse()
		fmt.Println("customJob2 returned : ", resp2)
		assert.Equal(t, input%2 == 0, resp2, "custom job 2 did not match")

		resultForCustomJob3 := results["customJob3"]
		resp3 := resultForCustomJob3.funcResponse()
		fmt.Println("customJob3 returned : ", resp3)
		assert.Equal(t, input%2 != 0, resp3, "custom job 3 did not match")
	}
}

func job1(val int) (func() interface{}, error) {
	fmt.Println("executing job1")
	time.Sleep(time.Second * 3)
	if val > 10 {
		return func() interface{} {
			return "success"
		}, nil
	}

	errMsg := fmt.Sprintf("too less for val in func1 : %v. It needs greater than 10 ", val)
	return nil, customErrorNew(errMsg, false)
}

func job2(val int) (func() interface{}, error) {
	fmt.Println("executing job2")
	time.Sleep(time.Second * 2)
	if val%2 == 0 {
		return func() interface{} {
			return "success"
		}, nil
	}
	errMsg := fmt.Sprintf("CRITICAL ERROR!! Val not divisible by 2 in func2 : %v", val)
	return nil, customErrorNew(errMsg, true)
}

func job3() (func() interface{}, error) {
	fmt.Println("executing job3")
	time.Sleep(time.Second * 2)
	return func() interface{} {
		return "func3 always passes!"
	}, nil
}

func customJob1(val int) (func() interface{}, error) {
	fmt.Println("executing customJob1")
	time.Sleep(time.Second * 3)
	localVal := 10 * val

	return func() interface{} {
		return localVal
	}, nil
}

func customJob2(val int) (func() interface{}, error) {
	fmt.Println("executing customJob2")
	time.Sleep(time.Second * 1)

	isValEven := false

	if val%2 == 0 {
		isValEven = true
	}

	return func() interface{} {
		return isValEven
	}, nil
}

func customJob3(val int) (func() interface{}, error) {
	fmt.Println("executing customJob3")
	time.Sleep(time.Second * 1)

	isValOdd := false

	if val%2 != 0 {
		isValOdd = true
	}

	return func() interface{} {
		return isValOdd
	}, nil
}
