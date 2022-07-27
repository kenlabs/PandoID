package peer

import (
	"fmt"
	"github.com/kenlabs/pando-id/pkg/did"
	pandoPID "github.com/kenlabs/pando/pkg/system"
	"github.com/multiformats/go-multibase"
	"regexp"
)

const (
	numAlgo   = "0"
	transform = multibase.Base58BTC
	didMethod = "peer"
)

func NewPeerDID() (didStr string, privateKey string, err error) {
	var peerID string

	peerID, privateKey, err = pandoPID.CreateIdentity()
	if err != nil {
		return
	}

	idStr := fmt.Sprint(numAlgo, string(transform), peerID)
	peerDID := &did.DID{ID: idStr, Method: didMethod}
	didStr = peerDID.String()

	return
}

func NewPeerDIDWithPeerID(peerID string) (didStr string, err error) {
	idStr := fmt.Sprint(numAlgo, string(transform), peerID)
	peerDID := &did.DID{ID: idStr, Method: didMethod}
	didStr = peerDID.String()

	return
}

func PeerDIDIsValid(peerDID string) (bool, error) {
	pattern := `^did:peer:(([01](z)([1-9a-km-zA-HJ-NP-Z]{46,47}))|(2((\.[AEVID](z)([1-9a-km-zA-HJ-NP-Z]{46,47}))+(\.(S)[0-9a-zA-Z=]*)?)))$
`
	match, err := regexp.Match(pattern, []byte(peerDID))
	if err != nil {
		return false, err
	}

	return match, nil
}
