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

//Run once before each test
func (suite *testCase) SetupTest() {
	//since jobNoError doesn't require an input, we have to wrap it around a function
	job3Wrapper := func(int) (func() interface{}, error) {
		return jobNoError()
	}

	barrier := NewBarrier().Add(jobGreater10WError).Add(jobEvenWErrorCritical).Add(job3Wrapper)
	suite.barrier = barrier
}

func (suite *testCase) TestExecute1() {

	barrier := suite.barrier
	t := suite.T()

	_, err := barrier.Execute(11)

	assert.NotNil(t, err, "Error should not be nil. ")
	assert.IsType(t, err, &customError{}, "Should be of custom error type")

	customErr := err.(*customError)
	assert.True(t, customErr.critical, "Error should be critical")
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

func (suite *testCaseCustom) TestCustom1() {
	t := suite.T()

	barrier := NewBarrier().AddN("jobMultiply10NoError", jobMultiply10NoError).
		AddN("jobEvenNoErrorWValue", jobEvenNoErrorWValue)
	jobOddNoErrorFunc := barrier.AddWNameReturned(jobOddNoError)

	input := 12
	results, err := barrier.Execute(input)
	assert.Nil(t, err, "err should be nil")

	if err != nil { //For completion sake
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
		resultForCustomJob1 := results["jobMultiply10NoError"]
		resp1 := resultForCustomJob1.funcResponse()
		assert.Equal(t, input*10, resp1, "custom job 1 did not match")

		resultForCustomJob2 := results["jobEvenNoErrorWValue"]
		resp2 := resultForCustomJob2.funcResponse()
		assert.Equal(t, input%2 == 0, resp2, "custom job 2 did not match")

		resultForCustomJob3 := results[jobOddNoErrorFunc]
		resp3 := resultForCustomJob3.funcResponse()
		assert.Equal(t, input%2 != 0, resp3, "custom job 3 did not match")
	}
}

func (suite *testCaseCustom) TestMix() {
	t := suite.T()

	barrier := NewBarrier().AddN("jobEvenWErrorCritical", jobEvenWErrorCritical).Add(jobEvenNoErrorWValue)
	barrier.AddWNameReturned(jobOddNoError)

	input := 11
	results, err := barrier.Execute(input)
	assert.NotNil(t, err, "err should not be nil")
	assert.Nil(t, results, "results should be nil since there was an error")
	assert.IsType(t, err, &customError{}, "Should be of custom error type")
	customErr := err.(*customError)
	assert.True(t, customErr.critical, "It should be critical")
}

func (suite *testCaseCustom) TestWithResults() {
	t := suite.T()

	barrier := NewBarrier().AddN("jobMultiply10NoError", jobMultiply10NoError).
		AddN("jobEvenWErrorCritical", jobEvenWErrorCritical).
		AddN("jobEvenNoErrorWValue", jobEvenNoErrorWValue)

	input := 11
	results := barrier.ExecuteAndReturnResults(input)

	for funcName, result := range results {
		switch funcName {
		case "jobEvenWErrorCritical":
			assert.NotNil(t, result.err, "should return an error")
		case "jobEvenNoErrorWValue":
			resp := result.funcResponse()
			respBool, ok := resp.(bool)
			assert.True(t, ok, "should be a boolean")
			assert.False(t, respBool, "response should be false")
		case "jobMultiply10NoError":
			assert.Nil(t, result.err, "should not return an error")
			assert.Equal(t, input*10, result.funcResponse(), "values do not match")
		}
	}
}

func (suite *testCaseCustom) TestWithNameReturn() {
	t := suite.T()

	barrier := NewBarrier()
	jobMultiply10NoErrorFunc := barrier.AddWNameReturned(jobMultiply10NoError)
	jobEvenWErrorCriticalFunc := barrier.AddWNameReturned(jobEvenWErrorCritical)
	jobEvenNoErrorWValueFunc := barrier.AddWNameReturned(jobEvenNoErrorWValue)

	input := 11
	results := barrier.ExecuteAndReturnResults(input)

	for funcName, result := range results {
		switch funcName {
		case jobEvenWErrorCriticalFunc:
			assert.NotNil(t, result.err, "should return an error")
		case jobEvenNoErrorWValueFunc:
			resp := result.funcResponse()
			respBool, ok := resp.(bool)
			assert.True(t, ok, "should be a boolean")
			assert.False(t, respBool, "response should be false")
		case jobMultiply10NoErrorFunc:
			assert.Nil(t, result.err, "should not return an error")
			assert.Equal(t, input*10, result.funcResponse(), "values do not match")
		}
	}
}

func jobGreater10WError(val int) (func() interface{}, error) {
	fmt.Println("executing jobGreater10WError")
	time.Sleep(time.Second * 3)
	if val > 10 {
		return func() interface{} {
			return "success"
		}, nil
	}

	errMsg := fmt.Sprintf("too less for val in func1 : %v. It needs greater than 10 ", val)
	return nil, customErrorNew(errMsg, false)
}

func jobEvenWErrorCritical(val int) (func() interface{}, error) {
	fmt.Println("executing jobEvenWErrorCritical")
	time.Sleep(time.Second * 2)
	if val%2 == 0 {
		return func() interface{} {
			return "success"
		}, nil
	}
	errMsg := fmt.Sprintf("CRITICAL ERROR!! Val not divisible by 2 in func2 : %v", val)
	return nil, customErrorNew(errMsg, true)
}

func jobNoError() (func() interface{}, error) {
	fmt.Println("executing jobNoError")
	time.Sleep(time.Second * 2)
	return func() interface{} {
		return "func3 always passes!"
	}, nil
}

func jobMultiply10NoError(val int) (func() interface{}, error) {
	fmt.Println("executing jobMultiply10NoError")
	time.Sleep(time.Second * 3)
	localVal := 10 * val

	return func() interface{} {
		return localVal
	}, nil
}

func jobEvenNoErrorWValue(val int) (func() interface{}, error) {
	fmt.Println("executing jobEvenNoErrorWValue")
	time.Sleep(time.Second * 1)

	isValEven := false

	if val%2 == 0 {
		isValEven = true
	}

	return func() interface{} {
		return isValEven
	}, nil
}

func jobOddNoError(val int) (func() interface{}, error) {
	fmt.Println("executing jobOddNoError")
	time.Sleep(time.Second * 1)

	isValOdd := false

	if val%2 != 0 {
		isValOdd = true
	}

	return func() interface{} {
		return isValOdd
	}, nil
}
