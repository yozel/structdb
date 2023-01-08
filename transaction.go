package structdb

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

// type ITransaction interface {
// 	Get(kn KindName) (Object, error)
// 	List(kn KindName) ([]Object, error)
// 	Set(obj Object) error
// 	Delete(kn KindName) error
// }

// var _ ITransaction = &Transaction{}

type Transaction struct {
	storage *Storage
	txn     *badger.Txn
}

func (txn *Transaction) Get(kn KindName) (Object, error) {
	if !txn.storage.Kinds.IsRegistered(kn.Kind) {
		return nil, fmt.Errorf("kind %s is not registered", kn.Kind)
	}
	it, err := txn.txn.Get(kn.key())
	if err != nil && err == badger.ErrKeyNotFound {
		log.Printf("key not found: %s", kn.key())
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return txn.storage.badgerItemToObject(it)
}

func (txn *Transaction) List(kn KindName) ([]Object, error) {
	objects := []Object{}
	if kn.Kind != "" && !txn.storage.Kinds.IsRegistered(kn.Kind) {
		return nil, fmt.Errorf("kind %s is not registered", kn.Kind)
	}

	if kn.Name != "" {
		if kn.Kind == "" {
			return nil, fmt.Errorf("kind is empty but name is not")
		}
		obj, err := txn.Get(kn)
		if err != nil {
			return nil, err
		}
		if obj != nil {
			objects = append(objects, obj)
		}
		return objects, nil
	} else {
		prefix := ""
		if kn.Kind != "" {
			prefix = string(kn.Kind) + "/"
		}
		key := []byte(prefix)
		it := txn.txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(key); it.ValidForPrefix(key); it.Next() {
			o, err := txn.storage.badgerItemToObject(it.Item())
			if err != nil {
				log.Printf("error getting object: %s", err)
				continue
			}
			objects = append(objects, o)
		}
		return objects, nil
	}
}

func (txn *Transaction) Set(obj Object) error {
	entry, err := txn.storage.objectToBadgerEntry(obj)
	if err != nil {
		return err
	}
	return txn.txn.SetEntry(entry)
}

func (txn *Transaction) Delete(kn KindName) error {
	if !txn.storage.Kinds.IsRegistered(kn.Kind) {
		return fmt.Errorf("kind %s is not registered", kn.Kind)
	}
	return txn.txn.Delete(kn.key())
}
