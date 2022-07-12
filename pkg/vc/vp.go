package vc

import (
	"encoding/json"

	"github.com/kenlabs/pando-id/pkg/internal/marshal"
	"github.com/kenlabs/pando-id/pkg/types"
)

const VerifiablePresentationType = "VerifiablePresentation"

func VerifiablePresentationTypeV1URI() types.URI {
	parsedURI, err := types.ParseURI(VerifiablePresentationType)
	if err != nil {
		panic(err)
	}
	return *parsedURI
}

type VerifiablePresentation struct {
	Context              []types.URI            `json:"@context"`
	ID                   *types.URI             `json:"id,omitempty"`
	Type                 []types.URI            `json:"type"`
	Holder               *types.URI             `json:"holder,omitempty"`
	VerifiableCredential []VerifiableCredential `json:"verifiableCredential,omitempty"`
	Proof                []interface{}          `json:"proof,omitempty"`
}

func (vp VerifiablePresentation) Proofs() ([]Proof, error) {
	var (
		target []Proof
		err    error
		asJSON []byte
	)
	asJSON, err = json.Marshal(vp.Proof)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(asJSON, &target)
	return target, err
}

func (vp VerifiablePresentation) MarshalJSON() ([]byte, error) {
	type alias VerifiablePresentation
	tmp := alias(vp)
	if data, err := json.Marshal(tmp); err != nil {
		return nil, err
	} else {
		return marshal.NormalizeDocument(data, pluralContext, marshal.Unplural(typeKey), marshal.Unplural(verifiableCredentialKey), marshal.Unplural(proofKey))
	}
}

func (vp *VerifiablePresentation) UnmarshalJSON(b []byte) error {
	type Alias VerifiablePresentation
	normalizedVC, err := marshal.NormalizeDocument(b, pluralContext, marshal.Plural(typeKey), marshal.Plural(verifiableCredentialKey), marshal.Plural(proofKey))
	if err != nil {
		return err
	}
	tmp := Alias{}
	err = json.Unmarshal(normalizedVC, &tmp)
	if err != nil {
		return err
	}
	*vp = (VerifiablePresentation)(tmp)
	return nil
}

func (vp VerifiablePresentation) UnmarshalProofValue(target interface{}) error {
	if asJSON, err := json.Marshal(vp.Proof); err != nil {
		return err
	} else {
		return json.Unmarshal(asJSON, target)
	}
}

func (vp VerifiablePresentation) IsType(vcType types.URI) bool {
	for _, t := range vp.Type {
		if t.String() == vcType.String() {
			return true
		}
	}

	return false
}

func (vp VerifiablePresentation) ContainsContext(context types.URI) bool {
	for _, c := range vp.Context {
		if c.String() == context.String() {
			return true
		}
	}

	return false
}
