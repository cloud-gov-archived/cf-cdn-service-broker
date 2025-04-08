package utils_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	. "github.com/cloud-gov/cf-cdn-service-broker/utils"
)

func TestHeaders(t *testing.T) {
	suite.Run(t, new(HeadersSuite))
}

type HeadersSuite struct {
	suite.Suite
}

func (h *HeadersSuite) SetupTest() {}

func (h *HeadersSuite) TestAdd() {
	headers := Headers{}
	headers.Add("abc-def")
	h.Equal(headers, Headers{"Abc-Def": true})
}

func (h *HeadersSuite) TestContains() {
	headers := Headers{"Abc-Def": true}
	h.True(headers.Contains("Abc-Def"))
	h.False(headers.Contains("Ghi-Jkl"))
}

func (h *HeadersSuite) TestStrings() {
	headers := Headers{"Abc-Def": true, "User-Agent": true}
	headerStrings := headers.Strings()
	h.Contains(headerStrings, "Abc-Def")
	h.Contains(headerStrings, "User-Agent")
	h.Equal(len(headerStrings), 2)
}
