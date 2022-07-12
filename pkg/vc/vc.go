package vc

import (
	"encoding/json"
	"github.com/kenlabs/pando-id/pkg/types"
	"net/url"
	"time"

	"github.com/kenlabs/pando-id/pkg/internal/marshal"
)

const VerifiableCredentialType = "VerifiableCredential"

func VerifiableCredentialTypeV1URI() types.URI {
	parsedURI, err := types.ParseURI(VerifiableCredentialType)
	if err != nil {
		panic(err)
	}
	return *parsedURI
}

const VCContextV1 = "https://www.w3.org/2018/credentials/v1"

func VCContextV1URI() types.URI {
	if pURI, err := types.ParseURI(VCContextV1); err != nil {
		panic(err)
	} else {
		return *pURI
	}
}

type VerifiableCredential struct {
	Context           []types.URI       `json:"@context"`
	ID                *types.URI        `json:"id,omitempty"`
	Type              []types.URI       `json:"type"`
	Issuer            types.URI         `json:"issuer"`
	IssuanceDate      time.Time         `json:"issuanceDate"`
	ExpirationDate    *time.Time        `json:"expirationDate,omitempty"`
	CredentialStatus  *CredentialStatus `json:"credentialStatus,omitempty"`
	CredentialSubject []interface{}     `json:"credentialSubject"`
	Proof             []interface{}     `json:"proof"`
}

type CredentialStatus struct {
	ID   url.URL `json:"id"`
	Type string  `json:"type"`
}

func (vc VerifiableCredential) Proofs() ([]Proof, error) {
	var (
		target []Proof
		err    error
		asJSON []byte
	)
	asJSON, err = json.Marshal(vc.Proof)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(asJSON, &target)
	return target, err
}

func (vc VerifiableCredential) MarshalJSON() ([]byte, error) {
	type alias VerifiableCredential
	tmp := alias(vc)
	if data, err := json.Marshal(tmp); err != nil {
		return nil, err
	} else {
		return marshal.NormalizeDocument(data, pluralContext, marshal.Unplural(typeKey), marshal.Unplural(credentialSubjectKey), marshal.Unplural(proofKey))
	}
}

func (vc *VerifiableCredential) UnmarshalJSON(b []byte) error {
	type Alias VerifiableCredential
	normalizedVC, err := marshal.NormalizeDocument(b, pluralContext, marshal.Plural(typeKey), marshal.Plural(credentialSubjectKey), marshal.Plural(proofKey))
	if err != nil {
		return err
	}
	tmp := Alias{}
	err = json.Unmarshal(normalizedVC, &tmp)
	if err != nil {
		return err
	}
	*vc = (VerifiableCredential)(tmp)
	return nil
}

func (vc VerifiableCredential) UnmarshalProofValue(target interface{}) error {
	if asJSON, err := json.Marshal(vc.Proof); err != nil {
		return err
	} else {
		return json.Unmarshal(asJSON, target)
	}
}

func (vc VerifiableCredential) UnmarshalCredentialSubject(target interface{}) error {
	if asJSON, err := json.Marshal(vc.CredentialSubject); err != nil {
		return err
	} else {
		return json.Unmarshal(asJSON, target)
	}
}

func (vc VerifiableCredential) IsType(vcType types.URI) bool {
	for _, t := range vc.Type {
		if t.String() == vcType.String() {
			return true
		}
	}

	return false
}

func (vc VerifiableCredential) ContainsContext(context types.URI) bool {
	for _, c := range vc.Context {
		if c.String() == context.String() {
			return true
		}
	}

	return false
}
