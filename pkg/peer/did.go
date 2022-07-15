package peer

import (
	"encoding/base64"
	"github.com/kenlabs/pando-id/pkg/did"
	pandoPID "github.com/kenlabs/pando/pkg/system"
	"github.com/libp2p/go-libp2p-core/crypto"
	"time"
)

const PeerDIDMethod = "peer"

// did:peer:peerID

type PeerDID struct {
	did.DID
}

type PeerDocument struct {
	did.Document
	// These properties should not be involved in DID doc
	Created int64 `json:"created,omitempty"`
	Updated int64 `json:"updated,omitempty"`
}

func NewPeerDID() (didStr string, peerID string, publicKey []byte, privateKey []byte, err error) {
	peerID, privateKeyStr, err := pandoPID.CreateIdentity()
	if err != nil {
		return
	}

	privateKeyEncoded, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return
	}

	privateKeyRaw, err := crypto.UnmarshalPrivateKey(privateKeyEncoded)
	privateKey, err = privateKeyRaw.Raw()
	if err != nil {
		return
	}
	publicKeyRaw := privateKeyRaw.GetPublic()
	publicKey, err = publicKeyRaw.Raw()
	if err != nil {
		return
	}

	peerDID := &PeerDID{}
	peerDID.Method = PeerDIDMethod
	peerDID.ID = peerID

	peerDocument := &PeerDocument{}
	peerDocument.Created = time.Now().UnixNano()
	peerDocument.ID = peerDID.DID
	peerDocument.Context = append(peerDocument.Context, did.DIDContextV1URI())

	return
}
