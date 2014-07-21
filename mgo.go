package dal

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type iter struct {
	Iter
	iter *mgo.Iter
}

func (i *iter) Next(inter interface{}) bool {
	return i.iter.Next(inter)
}

type query struct {
	Query
	query *mgo.Query
}

func (q *query) One(i interface{}) error {
	err := q.query.One(i)
	if err == mgo.ErrNotFound {
		err = ErrNotFound
	}
	return err
}

func (q *query) Iter() Iter {
	i := q.query.Iter()
	return &iter{iter: i}
}

func (q *query) Sort(s ...string) Query {
	q2 := q.query.Sort(s...)
	return &query{query: q2}
}

type collection struct {
	Collection
	col *mgo.Collection
}

func (c *collection) Find(b BSON) Query {
	q := c.col.Find(b)
	return &query{query: q}
}

type database struct {
	Database
	db *mgo.Database
}

func (d *database) C(name string) Collection {
	col := d.db.C(name)
	return &collection{col: col}
}

type session struct {
	Session
	mgoSession *mgo.Session
}

func (s *session) Clone() Session {
	a := s.mgoSession.Clone()
	return &session{mgoSession: a}
}

func (s *session) Close() {
	s.mgoSession.Close()
}

func (s *session) DB(name string) Database {
	db := s.mgoSession.DB(name)
	return &database{db: db}
}

type dal struct {
	DAL
}

func (d *dal) Connect(s string) (Session, error) {
	mgoSession, err := mgo.Dial(s)
	return &session{mgoSession: mgoSession}, err
}

func NewDAL() DAL {
	return &dal{}
}

type ObjectId string

func (id ObjectId) GetBSON() (interface{}, error) {
	return bson.ObjectId(id), nil
}

func (id ObjectId) Hex() string {
	return bson.ObjectId(id).Hex()
}

func (id ObjectId) Valid() bool {
	return bson.ObjectId(id).Valid()
}
