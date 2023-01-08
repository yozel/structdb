package structdb

import (
	"fmt"
	"strings"
)

type Kind string

type Name string

type KindName struct {
	Kind Kind
	Name Name
}

func NewKindName(kind, name string) KindName {
	return KindName{
		Kind: Kind(strings.ToLower(strings.TrimSpace(kind))),
		Name: Name(strings.ToLower(strings.TrimSpace(name))),
	}
}

func (kn KindName) key() []byte {
	return []byte(fmt.Sprintf("%s/%s", kn.Kind, kn.Name))
}
