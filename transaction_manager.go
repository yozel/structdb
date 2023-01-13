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
	kind, err := m.getKind()
	if err != nil {
		return nil, err
	}
	o, err := m.txn.Get(KindName{Kind: kind, Name: name})
	if err != nil {
		return
	}
	return objCast[T](o), nil
}

func (m TransactionManager[T]) List() ([]T, error) {
	kind, err := m.getKind()
	if err != nil {
		return nil, err
	}
	objs, err := m.txn.List(KindName{Kind: kind})
	if err != nil {
		return nil, err
	}
	return objCastSlice[T](objs)
}

func (m TransactionManager[T]) Set(obj T) error {
	kind, err := m.txn.storage.Kinds.ObjectToKind(obj)
	if err != nil {
		return err
	}
	kind2, err := m.getKind()
	if err != nil {
		return err
	}
	if kind != kind2 {
		return fmt.Errorf("object is not of type %s but %s: %+v", kind2, kind, obj)
	}
	return m.txn.Set(obj)
}

func (m TransactionManager[T]) Delete(name Name) error {
	kind, err := m.getKind()
	if err != nil {
		return err
	}
	return m.txn.Delete(KindName{Kind: kind, Name: name})
}

func (m TransactionManager[T]) getKind() (Kind, error) {
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
