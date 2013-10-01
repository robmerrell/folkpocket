package config

import (
	. "launchpad.net/gocheck"
	"os"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type configSuite struct{}

var _ = Suite(&configSuite{})

func (s *configSuite) SetUpSuite(c *C) {
	LoadConfigFile("testdata/config.toml")
}

func (s *configSuite) TestLoadConfigFile(c *C) {
	c.Check(tomlConfig.Get("test").(string), Equals, "string!")
}

func (s *configSuite) TestEnvValue(c *C) {
	os.Setenv("FOLKPOCKET_ENV", "testenv")

	c.Check(Env().Get("thekey").(string), Equals, "value")
}
