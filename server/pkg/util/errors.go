package util

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ErrTrace(err error) []string {
	if err == nil {
		return nil
	}
	trace := []string{}
	if errStack, ok := err.(stackTracer); ok {
		for _, f := range errStack.StackTrace() {
			trace = append(trace, strings.Split(fmt.Sprintf("%+v", f), "\n\t")...)
		}
	}
	return trace
}
