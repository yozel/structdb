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
