package peer

import (
	"fmt"
	"github.com/kenlabs/pando-id/pkg/did"
	"github.com/kenlabs/pando-id/pkg/test"
	"testing"
)

func Test_NewPeerDID(t *testing.T) {
	t.Run("pass if a generated peer did has proper prefix", func(t *testing.T) {
		didStr, prvKey, err := NewPeerDID()
		if err != nil {
			t.Fail()
		}
		fmt.Printf(
			"didStr:%s\nprvKey:%s\n",
			didStr, prvKey,
		)
		valid, err := DIDIsValid(didStr)
		if err != nil {
			t.Fail()
		}
		test.Assert(t, true, valid)

		d, err := did.Parse(didStr)
		if err != nil {
			t.Fail()
		}
		test.Assert(t, "peer", d.Method)
	})
}
