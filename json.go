package structdb

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func (s *Kinds) JsonToObject(j json.RawMessage) (obj Object, err error) {
	kind := gjson.GetBytes(j, "Kind").String()
	if kind == "" {
		return nil, fmt.Errorf("kind is empty")
	}
	name := gjson.GetBytes(j, "Metadata.Name").String()
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	return s.jsonToObject(NewKindName(kind, name), j)
}

func (s *Kinds) jsonToObject(kn KindName, j json.RawMessage) (obj Object, err error) {
	if !s.IsRegistered(kn.Kind) {
		return nil, fmt.Errorf("kind %s is not registered", kn.Kind)
	}
	j, err = sjson.SetBytes(j, "Metadata.Name", kn.Name)
	if err != nil {
		return nil, err
	}
	j, err = sjson.SetBytes(j, "Kind", kn.Kind)
	if err != nil {
		return nil, err
	}
	t, err := s.KindToType(kn.Kind)
	if err != nil {
		return nil, err
	}
	obj = reflect.New(t).Interface().(Object)
	err = json.Unmarshal(j, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
