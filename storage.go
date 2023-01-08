package structdb

import (
	"errors"
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
	k.mu.Lock()
	defer k.mu.Unlock()
	k.m[Kind(vv)] = v
	k.m2[reflect.TypeOf(v)] = Kind(vv)
	return nil
}

func (k *Kinds) IsRegistered(kind Kind) bool {
	k.mu.RLock()
	defer k.mu.RUnlock()
	_, ok := k.m[kind]
	return ok
}

func (k *Kinds) KindToType(kind Kind) reflect.Type {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return reflect.TypeOf(k.m[kind])
}

func (k *Kinds) ObjectToKind(v Object) Kind {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.m2[reflect.TypeOf(v)]
}

type Storage struct {
	db    *badger.DB
	Kinds Kinds
}

func New(path string) (*Storage, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
		Kinds: Kinds{
			m:  make(map[Kind]Object),
			m2: make(map[reflect.Type]Kind),
		},
	}, nil
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
