package did

import "fmt"

type Param struct {
	Name  string
	Value string
}

func (p *Param) String() string {
	if p.Name == "" {
		return ""
	}

	if 0 < len(p.Value) {
		return fmt.Sprint(p.Name, "+", p.Value)
	}

	return p.Name
}
