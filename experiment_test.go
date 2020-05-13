package exp_test

import (
	exp "experiment"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestGobalInitIsWorkInUnitTest(t *testing.T) {
	suite.Run(t, new(GoobalInitIsWorkInUnitTestTestSuite))
}

type GoobalInitIsWorkInUnitTestTestSuite struct {
	suite.Suite
}

func (ts *GoobalInitIsWorkInUnitTestTestSuite) SetupTest() {

}

func (ts *GoobalInitIsWorkInUnitTestTestSuite) Test_Scenario() {
	expected := 20
	actual := exp.Golbal["aaa"]
	ts.Assert().Equal(expected, actual)
}
