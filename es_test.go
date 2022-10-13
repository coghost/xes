package xes_test

import (
	"testing"

	"github.com/coghost/xes"
	"github.com/stretchr/testify/suite"
)

type XesSuite struct {
	suite.Suite
}

func TestXes(t *testing.T) {
	suite.Run(t, new(XesSuite))
}

func (s *XesSuite) SetupSuite() {
}

func (s *XesSuite) TearDownSuite() {
}

func (s *XesSuite) Test_00_init() {
	es, err := xes.NewEsLogger()
	s.Nil(es)
	s.NotNil(err)
	s.Equal(xes.ErrServersRequired, err)

	es, err = xes.NewEsLogger(xes.WithUrls(""))
	s.Nil(es)
	s.NotNil(err)
	s.Equal(xes.ErrServerIsEmpty, err)
}
