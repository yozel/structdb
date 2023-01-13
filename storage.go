package structdb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger/v3"
)

// type IStorage interface {
// 	TxnRO(fn func(txn *Transaction) error) error
// 	TxnRW(fn func(txn *Transaction) error) error
// 	Close()
// }

// var _ IStorage = &Storage{}

type Kinds struct {
	mu sync.RWMutex
	m  map[Kind]Object
	m2 map[reflect.Type]Kind
}

var ErrInvalidKind = errors.New("invalid kind")

func (k *Kinds) Register(vv string, v Object) error {
	if vv != strings.TrimSpace(strings.ToLower(vv)) {
		return ErrInvalidKind
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return fmt.Errorf("Kind %s is a pointer, must be a struct", v)
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.m[Kind(vv)] = v
	k.m2[reflect.TypeOf(v)] = Kind(vv)
	return nil
}

func (k *Kinds) isRegistered(kind Kind) bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	_, ok := k.m[kind]
	return ok
}

func (k *Kinds) NewObject(kind Kind) (Object, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	v, ok := k.m[kind]
	if !ok {
		return nil, ErrInvalidKind
	}
	return reflect.New(reflect.TypeOf(v)).Interface().(Object), nil
}

func (k *Kinds) GetKind(v Object) (Kind, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	kind, ok := k.m2[t]
	if !ok {
		return "", ErrInvalidKind
	}
	return kind, nil
}

type Storage struct {
	db    *badger.DB
	Kinds Kinds
}

func New(db *badger.DB) *Storage {
	return &Storage{
		db: db,
		Kinds: Kinds{
			m:  make(map[Kind]Object),
			m2: make(map[reflect.Type]Kind),
		},
	}
}

func (s *Storage) TxnRO(fn func(*Transaction) error) error {
	return s.db.View(func(txn *badger.Txn) error {
		return fn(&Transaction{
			storage: s,
			txn:     txn,
		})
	})
}

func (s *Storage) TxnRW(fn func(*Transaction) error) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return fn(&Transaction{
			storage: s,
			txn:     txn,
		})
	})
}

func (s *Storage) Get(kn KindName) (Object, error) {
	var obj Object
	err := s.TxnRO(func(txn *Transaction) error {
		var err error
		obj, err = txn.Get(kn)
		return err
	})
	return obj, err
}

func (s *Storage) List(kn KindName) ([]Object, error) {
	var objs []Object
	err := s.TxnRO(func(txn *Transaction) error {
		var err error
		objs, err = txn.List(kn)
		return err
	})
	return objs, err
}

func (s *Storage) GetKindName(obj Object) (*KindName, error) {
	kind, err := s.Kinds.GetKind(obj)
	if err != nil {
		return nil, err
	}
	return &KindName{
		Kind: kind,
		Name: obj.GetMetadata().Name,
	}, nil
}

func (s *Storage) Set(obj Object) error {
	return s.TxnRW(func(txn *Transaction) error {
		return txn.Set(obj)
	})
}

func (s *Storage) Delete(kn KindName) error {
	return s.TxnRW(func(txn *Transaction) error {
		return txn.Delete(kn)
	})
}

func (s *Storage) Close() {
	s.db.Close()
}
