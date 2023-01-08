package structdb

import (
	"time"
)

type ObjectType struct {
	Kind        Kind `json:"Kind" yaml:"Kind"`
	*ObjectMeta `json:"Metadata" yaml:"Metadata"`
}

type ObjectMeta struct {
	Name      Name      `json:"Name" yaml:"Name"`
	CreatedAt time.Time `json:"CreatedAt" yaml:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" yaml:"UpdatedAt"`
}

func (o *ObjectMeta) GetMetadata() *ObjectMeta {
	return o
}

type Object interface {
	GetMetadata() *ObjectMeta
}
