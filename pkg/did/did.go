package did

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kenlabs/pando-id/pkg/types"
	"net/url"
	"strings"
)

const DIDContextV1 = "https://www.w3.org/ns/did/v1"

func DIDContextV1URI() types.URI {
	parsedURI, err := types.ParseURI(DIDContextV1)
	if err != nil {
		panic(err)
	}
	return *parsedURI
}

type DID struct {
	Method       string
	ID           string
	IDStrings    []string
	Params       []Param
	Path         string
	PathSegments []string
	Query        string
	Fragment     string
}

func (d *DID) String() string {
	var buf strings.Builder
	buf.WriteString("did:")

	if d.Method != "" {
		buf.WriteString(d.Method)
		buf.WriteByte(':')
	} else {
		return ""
	}

	if d.ID != "" {
		buf.WriteString(d.ID)
	} else if len(d.IDStrings) > 0 {
		buf.WriteString(strings.Join(d.IDStrings[:], ":"))
	} else {
		return ""
	}

	if len(d.Params) > 0 {
		for _, p := range d.Params {
			param := p.String()
			if param != "" {
				buf.WriteByte(';')
				buf.WriteString(param)
			} else {
				return ""
			}
		}
	}

	if d.Path != "" {
		buf.WriteByte('/')
		buf.WriteString(d.Path)
	} else if len(d.PathSegments) > 0 {
		buf.WriteByte('/')
		buf.WriteString(strings.Join(d.PathSegments[:], "/"))
	}

	if d.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(d.Query)
	}

	if d.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(d.Fragment)
	}

	return buf.String()
}

func Parse(input string) (*DID, error) {
	p := &parser{
		input: input,
		out:   &DID{},
	}

	parserState := p.checkLength
	for parserState != nil {
		parserState = parserState()
	}

	err := p.err
	if err != nil {
		return nil, err
	}

	p.out.ID = strings.Join(p.out.IDStrings[:], ":")
	p.out.Path = strings.Join(p.out.PathSegments[:], "/")

	if p.out.IsURL() {
		return nil, ErrInvalidDID.wrap(errors.New("DID cannot have path, fragment or query params"))
	}

	return p.out, nil
}

func (d *DID) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *DID) Equals(other DID) bool {
	return d.String() == other.String()
}

func (d *DID) UnmarshalJSON(bytes []byte) error {
	var didString string
	err := json.Unmarshal(bytes, &didString)
	if err != nil {
		return ErrInvalidDID.wrap(err)
	}
	tmp, err := Parse(didString)
	if err != nil {
		return ErrInvalidDID.wrap(err)
	}
	*d = *tmp
	return nil
}

func (d *DID) MarshalJSON() ([]byte, error) {
	didAsString := d.String()
	return json.Marshal(didAsString)
}

func (d *DID) URI() types.URI {
	return types.URI{
		URL: url.URL{
			Scheme:   "did",
			Opaque:   fmt.Sprintf("%s:%s", d.Method, d.ID),
			Fragment: d.Fragment,
		},
	}
}

func (d *DID) Empty() bool {
	return d.Method == ""
}

func (d *DID) IsURL() bool {
	return len(d.Params) > 0 || d.Path != "" || len(d.PathSegments) > 0 || d.Query != "" || d.Fragment != ""
}
