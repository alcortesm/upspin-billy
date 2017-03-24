package upspinfs

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v2/test"
)

func Test(t *testing.T) { TestingT(t) }

type UpspinSuite struct {
	test.FilesystemSuite
}

var _ = Suite(&UpspinSuite{})

func (s *UpspinSuite) SetUpTest(c *C) {
	s.FilesystemSuite.Fs = New()
}

func (s *UpspinSuite) TearDownTest(c *C) {
	// do nothing
}
