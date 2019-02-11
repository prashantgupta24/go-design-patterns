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
	barrier := &Barrier{}

	//since job3 doesn't require an input, we have to wrap it around a function
	job3Wrapper := func(int) (func() interface{}, error) {
		return job3()
	}

	barrier.Add(job1).Add(job2).Add(job3Wrapper)
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
	barrier := &Barrier{}
	//t := suite.T()
	barrier.Add(customJob1).Add(customJob2).Add(customJob3)

	results := barrier.ExecuteAndReturnResults(12)

	//hasError := false
	for _, result := range results {
		if result.err != nil {
			//hasError = true
			if err, ok := result.err.(*customError); ok {
				if err.critical {
					fmt.Println("CRITICAL ERROR!! ", err)
				} else {
					fmt.Println("ERROR!! ", err)
				}
			} else {
				fmt.Println("ERROR!! >>> ", result.err)
			}
		} else {
			resp := result.response()
			switch valType := resp.(type) {
			case int:
				fmt.Println("function returned int : ", valType)
			case bool:
				fmt.Println("function returned bool : ", valType)
			}
		}
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
