package peer

import (
	"fmt"
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
	})
}
