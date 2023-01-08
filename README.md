# StructDB

StructDB is a Go library that provides a simple, unified API for storing and retrieving Golang structs in BadgerDB, a popular key-value store. It abstracts away the complexities of working with BadgerDB, allowing developers to focus on building their applications. With StructDB, you can easily store and retrieve structs from BadgerDB with just a few lines of code. StructDB is designed to be easy to use and highly performant, making it the ideal choice for Go developers looking to simplify their storage needs with BadgerDB.

Here's a quick example of how StructDB works:

```golang
package main

import (
    "fmt"
    "github.com/structdb/structdb"
)


type Storage struct {
	*storage.Storage
}

type User struct {
    storage.ObjectType
    ID   int64
    Name string
}

func (s *Storage) User() *storage.StorageManager[User] {
	return storage.NewStorageManager[User](s.Storage)
}

func NewStorage(s *storage.Storage) *Storage {
	ss := &Storage{s}
	ss.Kinds.Register(User{})
    return ss
}

func main() {
    s := storage.NewStorage()
    ss := NewStorage(s)

    ss.User().Set(&User{
        ID:   1,
        Name: "John",
    })
    ss.TxnRW(func(txn *Transaction) error {
        u := &User{
            ID:   1,
            Name: "John",
        }
        err := txn.Put(u)
        if err != nil {
            return err
        }
        return nil
    })
    ss.TxnRW(func(txn *Transaction) error {
        u := &User{}
        err := txn.Get(u, 1)
        if err != nil {
            return err
        }
        fmt.Println(u)
        return nil
    })
}
```
