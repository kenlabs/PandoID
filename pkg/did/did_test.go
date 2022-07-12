package did

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestDID_Parse(t *testing.T) {
	t.Run("pass if a did has proper prefix and length is greater than 7", func(t *testing.T) {
		d, err := Parse("did:a:a")
		assert(t, nil, err)
		assert(t, true, d != nil)
	})
}

func assert(t *testing.T, expected interface{}, actual interface{}, args ...interface{}) {
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
