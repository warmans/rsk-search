package flag

import (
	goflag "flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/util"
	"os"
	"strconv"
	"strings"
)

func StringVarEnv(flagsSet *pflag.FlagSet, s *string, prefix string, name string, value string, usage string) {
	flagsSet.StringVar(s, name, value, usage)
	stringFromEnv(s, prefix, name)
}

func BoolVarEnv(flagsSet *pflag.FlagSet, s *bool, prefix string, name string, value bool, usage string) {
	flagsSet.BoolVar(s, name, value, usage)
	boolFromEnv(s, prefix, name)
}

func Int64VarEnv(flagsSet *pflag.FlagSet, s *int64, prefix string, name string, value int64, usage string) {
	flagsSet.Int64Var(s, name, value, usage)
	int64FromEnv(s, prefix, name)
}

func stringFromEnv(p *string, prefix, name string) {
	if prefix != "" {
		prefix = "_" + strings.ToUpper(prefix)
	}
	val := os.Getenv(fmt.Sprintf("%s%s", prefix, strings.ToUpper(strings.Replace(name, "-", "_", -1))))
	if val == "" {
		return
	}
	valPtr := &val
	*p = *valPtr
}

func boolFromEnv(p *bool, prefix, name string) {
	if prefix != "" {
		prefix = "_" + strings.ToUpper(prefix)
	}
	val := os.Getenv(fmt.Sprintf("%s%s", prefix, strings.ToUpper(strings.Replace(name, "-", "_", -1))))
	if val == "" {
		return
	}
	boolVal := val == "true"
	valPtr := &boolVal
	*p = *valPtr
}

func int64FromEnv(p *int64, prefix, name string) {
	if prefix != "" {
		prefix = "_" + strings.ToUpper(prefix)
	}
	val := os.Getenv(fmt.Sprintf("%s%s", prefix, strings.ToUpper(strings.Replace(name, "-", "_", -1))))
	if val == "" {
		return
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return
	}
	valPtr := util.Int64P(int64(intVal))
	*p = *valPtr
}

func Parse() {
	goflag.Parse()
}
