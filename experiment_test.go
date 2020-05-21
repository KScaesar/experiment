package main_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestReadFile(t *testing.T) {
	suite.Run(t, new(GoobalInitIsWorkInUnitTestTestSuite))
}

type GoobalInitIsWorkInUnitTestTestSuite struct {
	suite.Suite
}

func (ts *GoobalInitIsWorkInUnitTestTestSuite) Test_Scenario() {
}
