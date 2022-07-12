package did

import (
	"time"
)

type constError string

func (err constError) Error() string {
	return string(err)
}

const (
	InvalidDIDErr  = constError("supplied DID is invalid")
	NotFoundErr    = constError("supplied DID wasn't found")
	DeactivatedErr = constError("supplied DID is deactivated")
)

type Resolver interface {
	Resolve(inputDID string) (*Document, *DocumentMetadata, error)
}

type DocumentMetadata struct {
	Created    *time.Time
	Updated    *time.Time
	Properties map[string]interface{}
}
