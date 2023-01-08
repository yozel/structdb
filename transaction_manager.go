package structdb

import (
	"fmt"
	"reflect"
)

type TransactionManager[T Object] struct {
	txn *Transaction
}

func NewTransactionManager[T Object](txn *Transaction) *TransactionManager[T] {
	return &TransactionManager[T]{txn: txn}
}

func (m *TransactionManager[T]) Get(name Name) (r *T, err error) {
	o, err := m.txn.Get(KindName{Kind: m.getKind(), Name: name})
	if err != nil {
		return
	}
	return objCast[T](o), nil
}

func (m TransactionManager[T]) List() ([]T, error) {
	objs, err := m.txn.List(KindName{Kind: m.getKind()})
	if err != nil {
		return nil, err
	}
	return objCastSlice[T](objs)
}

func (m TransactionManager[T]) Set(obj T) error {
	if m.txn.storage.Kinds.ObjectToKind(obj) != m.getKind() {
		return fmt.Errorf("object is not of type %s but %s: %+v", m.getKind(), m.txn.storage.Kinds.ObjectToKind(obj), obj)
	}
	return m.txn.Set(obj)
}

func (m TransactionManager[T]) Delete(name Name) error {
	return m.txn.Delete(KindName{Kind: m.getKind(), Name: name})
}

func (m TransactionManager[T]) getKind() Kind {
	return m.txn.storage.Kinds.ObjectToKind((*new(T)))
}

func objCast[T any](o Object) *T {
	if o == nil {
		return nil
	}
	return reflect.ValueOf(o).Interface().(*T)
}

func objCastSlice[T any](o []Object) ([]T, error) {
	objects := []T{}
	for _, o := range o {
		obj := objCast[T](o)
		if obj == nil {
			return nil, fmt.Errorf("object is nil: %+v", o)
		}
		objects = append(objects, *obj)
	}
	return objects, nil
}
