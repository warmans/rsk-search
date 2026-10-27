package flag

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestBoolVarEnv(t *testing.T) {
	cases := []struct {
		name     string
		envName  string
		envValue string
		def      bool
		expected bool
	}{
		{"env true", "RSK_SEARCH_MY_BOOL", "true", false, true},
		{"env false", "RSK_SEARCH_MY_BOOL", "false", true, false},
		{"env unset keeps default", "", "", true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.envName != "" {
				t.Setenv(c.envName, c.envValue)
			}
			v := c.def
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			BoolVarEnv(fs, &v, "RSK_SEARCH", "my-bool", c.def, "usage")
			if v != c.expected {
				t.Errorf("got %v, expected %v", v, c.expected)
			}
		})
	}
}

func TestInt64VarEnv(t *testing.T) {
	cases := []struct {
		name     string
		envName  string
		envValue string
		def      int64
		expected int64
	}{
		{"env value", "RSK_SEARCH_MY_INT", "42", 0, 42},
		{"env invalid keeps default", "RSK_SEARCH_MY_INT", "not-a-number", 7, 7},
		{"env unset keeps default", "", "", 9, 9},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.envName != "" {
				t.Setenv(c.envName, c.envValue)
			}
			v := c.def
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			Int64VarEnv(fs, &v, "RSK_SEARCH", "my-int", c.def, "usage")
			if v != c.expected {
				t.Errorf("got %d, expected %d", v, c.expected)
			}
		})
	}
}

func TestStringVarEnv(t *testing.T) {
	t.Setenv("RSK_SEARCH_MY_STRING", "from-env")

	v := "default"
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	StringVarEnv(fs, &v, "RSK_SEARCH", "my-string", "default", "usage")
	if v != "from-env" {
		t.Errorf("got %s, expected from-env", v)
	}
}
