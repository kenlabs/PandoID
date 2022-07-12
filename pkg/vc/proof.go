package vc

import (
	"time"

	"github.com/kenlabs/pando-id/pkg/types"
)

type Proof struct {
	Type               types.ProofType `json:"type"`
	ProofPurpose       string          `json:"proofPurpose"`
	VerificationMethod types.URI       `json:"verificationMethod"`
	Created            time.Time       `json:"created"`
	Domain             *string         `json:"domain,omitempty"`
}

type JSONWebSignature2020Proof struct {
	Proof
	Jws string `json:"jws"`
}
