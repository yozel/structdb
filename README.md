# StructDB

StructDB is a Go library that provides a simple, unified API for storing and retrieving Golang structs in BadgerDB, a popular key-value store. It abstracts away the complexities of working with BadgerDB, allowing developers to focus on building their applications. With StructDB, you can easily store and retrieve structs from BadgerDB with just a few lines of code. StructDB is designed to be easy to use and highly performant, making it the ideal choice for Go developers looking to simplify their storage needs with BadgerDB.

Here's a quick example of how StructDB works:

```golang
package main

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/yozel/structdb"
)

type Storage struct {
	*structdb.Storage
}

type User struct {
	structdb.ObjectType
	Firstname string
	Lastname  string
}

func (s *Storage) User() *structdb.StorageManager[User] {
	return structdb.NewStorageManager[User](s.Storage)
}

func NewStorage(db *badger.DB) (*Storage, error) {
	ss := &Storage{
		Storage: structdb.New(db),
	}
	err := ss.Kinds.Register("user", User{})
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("badger.db"))
	if err != nil {
		panic(err)
	}
	ss, err := NewStorage(db)
	if err != nil {
		panic(err)
	}
	defer ss.Close()

	err = ss.User().Set(User{
		ObjectType: structdb.ObjectType{
			ObjectMeta: &structdb.ObjectMeta{
				Name: "jdoe",
			},
		},
		Firstname: "John",
		Lastname:  "Doe",
	})
	if err != nil {
		panic(err)
	}

	u, err := ss.User().Get("jdoe")
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Firstname, u.Lastname)

	err = ss.User().Delete("jdoe")
	if err != nil {
		panic(err)
	}
}
```
