package ssh

import (
	"fmt"
	"os"

	. "gopkg.in/check.v1"
)

type SuiteCommon struct{}

var _ = Suite(&SuiteCommon{})

func (s *SuiteCommon) TestKeyboardInteractiveName(c *C) {
	a := &KeyboardInteractive{
		User:      "test",
		Challenge: nil,
	}
	c.Assert(a.Name(), Equals, KeyboardInteractiveName)
}

func (s *SuiteCommon) TestKeyboardInteractiveString(c *C) {
	a := &KeyboardInteractive{
		User:      "test",
		Challenge: nil,
	}
	c.Assert(a.String(), Equals, fmt.Sprintf("user: test, name: %s", KeyboardInteractiveName))
}

func (s *SuiteCommon) TestPasswordName(c *C) {
	a := &Password{
		User: "test",
		Pass: "",
	}
	c.Assert(a.Name(), Equals, PasswordName)
}

func (s *SuiteCommon) TestPasswordString(c *C) {
	a := &Password{
		User: "test",
		Pass: "",
	}
	c.Assert(a.String(), Equals, fmt.Sprintf("user: test, name: %s", PasswordName))
}

func (s *SuiteCommon) TestPasswordCallbackName(c *C) {
	a := &PasswordCallback{
		User:     "test",
		Callback: nil,
	}
	c.Assert(a.Name(), Equals, PasswordCallbackName)
}

func (s *SuiteCommon) TestPasswordCallbackString(c *C) {
	a := &PasswordCallback{
		User:     "test",
		Callback: nil,
	}
	c.Assert(a.String(), Equals, fmt.Sprintf("user: test, name: %s", PasswordCallbackName))
}

func (s *SuiteCommon) TestPublicKeysName(c *C) {
	a := &PublicKeys{
		User:   "test",
		Signer: nil,
	}
	c.Assert(a.Name(), Equals, PublicKeysName)
}

func (s *SuiteCommon) TestPublicKeysString(c *C) {
	a := &PublicKeys{
		User:   "test",
		Signer: nil,
	}
	c.Assert(a.String(), Equals, fmt.Sprintf("user: test, name: %s", PublicKeysName))
}

func (s *SuiteCommon) TestPublicKeysCallbackName(c *C) {
	a := &PublicKeysCallback{
		User:     "test",
		Callback: nil,
	}
	c.Assert(a.Name(), Equals, PublicKeysCallbackName)
}

func (s *SuiteCommon) TestPublicKeysCallbackString(c *C) {
	a := &PublicKeysCallback{
		User:     "test",
		Callback: nil,
	}
	c.Assert(a.String(), Equals, fmt.Sprintf("user: test, name: %s", PublicKeysCallbackName))
}
func (s *SuiteCommon) TestNewSSHAgentAuth(c *C) {
	addr := os.Getenv("SSH_AUTH_SOCK")
	err := os.Unsetenv("SSH_AUTH_SOCK")
	c.Assert(err, IsNil)

	defer func() {
		err := os.Setenv("SSH_AUTH_SOCK", addr)
		c.Assert(err, IsNil)
	}()

	k, err := NewSSHAgentAuth("foo")
	c.Assert(k, IsNil)
	c.Assert(err, Equals, ErrEmptySSHAgentAddr)
}
