package did

import (
	"github.com/kenlabs/pando-id/pkg/test"
	"testing"
)

func TestDID_Parse(t *testing.T) {
	t.Run("pass if a did has proper prefix and length is greater than 7", func(t *testing.T) {
		d, err := Parse("did:a:a")
		test.Assert(t, nil, err)
		test.Assert(t, true, d != nil)
	})
}
