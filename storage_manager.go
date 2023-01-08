package structdb

type StorageManager[T Object] struct {
	s *Storage
}

func NewStorageManager[T Object](s *Storage) *StorageManager[T] {
	return &StorageManager[T]{s: s}
}

func (sm *StorageManager[T]) Get(name Name) (r *T, err error) {
	err = sm.s.TxnRO(func(txn *Transaction) (err error) {
		r, err = (&TransactionManager[T]{txn: txn}).Get(name)
		return
	})
	return r, err
}

func (sm *StorageManager[T]) List() (r []T, err error) {
	err = sm.s.TxnRO(func(txn *Transaction) (err error) {
		x := &TransactionManager[T]{txn: txn}
		r, err = x.List()
		return
	})
	return
}

func (sm *StorageManager[T]) Delete(name Name) (err error) {
	err = sm.s.TxnRW(func(txn *Transaction) (err error) {
		x := &TransactionManager[T]{txn: txn}
		err = x.Delete(name)
		return
	})
	return
}

func (sm *StorageManager[T]) Set(obj T) (err error) {
	err = sm.s.TxnRW(func(txn *Transaction) (err error) {
		x := &TransactionManager[T]{txn: txn}
		err = x.Set(obj)
		return
	})
	return
}
