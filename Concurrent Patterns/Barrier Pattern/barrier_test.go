package main

import (
	"testing"

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

	job3Wrapper := func(int) (string, error) {
		return job3()
	}

	barrier.add(job1).add(job2).add(job3Wrapper)
	suite.barrier = barrier
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
