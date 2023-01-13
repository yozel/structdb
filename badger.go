package structdb

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/tidwall/sjson"
)

var isAlphaNumeric = regexp.MustCompile(`^[a-z0-9-_]*$`)

func (s *Storage) badgerItemToObject(item *badger.Item) (Object, error) {
	parts := strings.Split(string(item.Key()), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("key is not a valid key")
	}
	kn := NewKindName(parts[0], parts[1])
	var obj Object
	err := item.Value(func(val []byte) error {
		var err error
		obj, err = s.Kinds.jsonToObject(kn, val)
		return err
	})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *Storage) objectToBadgerEntry(obj Object) (*badger.Entry, error) {
	kind, err := s.Kinds.ObjectToKind(obj)
	if err != nil {
		return nil, err
	}
	if !isAlphaNumeric.MatchString(string(obj.GetMetadata().Name)) {
		return nil, fmt.Errorf("name must be alphanumeric")
	}
	if obj.GetMetadata().CreatedAt.IsZero() {
		obj.GetMetadata().CreatedAt = time.Now().UTC()
	}
	obj.GetMetadata().UpdatedAt = time.Now().UTC()
	name := obj.GetMetadata().Name
	val, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	val, err = sjson.DeleteBytes(val, "Metadata.Name")
	if err != nil {
		return nil, err
	}
	val, err = sjson.DeleteBytes(val, "Kind")
	if err != nil {
		return nil, err
	}
	return badger.NewEntry(KindName{Kind: kind, Name: name}.key(), val), nil
}
