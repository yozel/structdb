package main

import (
	"fmt"

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

func NewStorage(s *structdb.Storage) *Storage {
	ss := &Storage{s}
	err := ss.Kinds.Register("user", User{})
	if err != nil {
		panic(err)
	}
	return ss
}

func main() {
	s, err := structdb.New("badger.db")
	if err != nil {
		panic(err)
	}
	defer s.Close()
	ss := NewStorage(s)

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
