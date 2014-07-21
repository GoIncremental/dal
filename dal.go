package dal

import (
	"errors"
)

var ErrNotFound = errors.New("ErrNotFound")

type BSON map[string]interface{}

type DAL interface {
	Connect(string) (Session, error)
}

type Session interface {
	Clone() Session
	Close()
	DB(s string) Database
}

type Database interface {
	C(string) Collection
}

type Collection interface {
	Find(BSON) Query
}

type Query interface {
	One(interface{}) error
	Sort(...string) Query
	Iter() Iter
}

type Iter interface {
	Next(interface{}) bool
}
