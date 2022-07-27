package test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Assert(t *testing.T, expected interface{}, actual interface{}, args ...interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		argsLength := len(args)
		var message string

		if argsLength == 1 {
			message = args[0].(string)
		}

		if argsLength > 1 {
			message = fmt.Sprintf(args[0].(string), args[1:]...)
		}

		if message != "" {
			message = "\t" + message + "\n\n"
		}

		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\tExpected: %#v\n\tActual: %#v\n%s", filepath.Base(file), line, expected, actual, message)
		t.FailNow()
	}
}
