package git

import (
	"github.com/src-d/go-git-fixtures"
	"srcd.works/go-git.v4/plumbing/transport"
	"srcd.works/go-git.v4/plumbing/transport/test"

	. "gopkg.in/check.v1"
)

type UploadPackSuite struct {
	test.UploadPackSuite
	fixtures.Suite
}

var _ = Suite(&UploadPackSuite{})

func (s *UploadPackSuite) SetUpSuite(c *C) {
	s.Suite.SetUpSuite(c)

	s.UploadPackSuite.Client = DefaultClient

	ep, err := transport.NewEndpoint("git://github.com/git-fixtures/basic.git")
	c.Assert(err, IsNil)
	s.UploadPackSuite.Endpoint = ep

	ep, err = transport.NewEndpoint("git://github.com/git-fixtures/empty.git")
	c.Assert(err, IsNil)
	s.UploadPackSuite.EmptyEndpoint = ep

	ep, err = transport.NewEndpoint("git://github.com/git-fixtures/non-existent.git")
	c.Assert(err, IsNil)
	s.UploadPackSuite.NonExistentEndpoint = ep

}
