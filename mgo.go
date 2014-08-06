package dal

import (
	"encoding/gob"
	"log"

	"github.com/goincremental/dal/Godeps/_workspace/src/labix.org/v2/mgo"
	"github.com/goincremental/dal/Godeps/_workspace/src/labix.org/v2/mgo/bson"
)

func init() {
	gob.Register(ObjectID(""))
}

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
	col *mgo.Collection
}

func (c *collection) Find(b BSON) Query {
	q := c.col.Find(b)
	return &query{query: q}
}

func (c *collection) EnsureIndex(index Index) error {
	i := mgo.Index{
		Key:         index.Key,
		Background:  index.Background,
		Sparse:      index.Sparse,
		ExpireAfter: index.ExpireAfter,
	}
	return c.col.EnsureIndex(i)
}

func (c *collection) FindID(id interface{}) Query {
	q := c.col.FindId(id)
	return &query{query: q}
}

func (c *collection) RemoveID(id interface{}) error {
	return c.col.RemoveId(id)
}

func (c *collection) UpsertID(id interface{}, update interface{}) (*ChangeInfo, error) {
	log.Printf("upsertId")
	mci, err := c.col.UpsertId(id, update)
	if err != nil {
		log.Printf("error upserting %s\n", err)
	}
	ci := &ChangeInfo{}
	if mci != nil {
		ci.Updated = mci.Updated
		ci.Removed = mci.Removed
		ci.UpsertedId = mci.UpsertedId
	}
	log.Printf("change info %s", ci)
	return ci, err
}

type database struct {
	Database
	db *mgo.Database
}

func (d *database) C(name string) Collection {
	col := d.db.C(name)
	return &collection{col: col}
}

type connection struct {
	Connection
	mgoSession *mgo.Session
}

func (c *connection) Clone() Connection {
	a := c.mgoSession.Clone()
	return &connection{mgoSession: a}
}

func (c *connection) Close() {
	c.mgoSession.Close()
}

func (c *connection) DB(name string) Database {
	db := c.mgoSession.DB(name)
	return &database{db: db}
}

type dal struct {
	DAL
}

func (d *dal) Connect(s string) (Connection, error) {
	log.Printf("Connect: %s\n", s)
	mgoSession, err := mgo.Dial(s)
	return &connection{mgoSession: mgoSession}, err
}

func NewDAL() DAL {
	return &dal{}
}

type ObjectID string

func (id ObjectID) GetBSON() (interface{}, error) {
	return bson.ObjectId(id), nil
}

func (id ObjectID) Hex() string {
	return bson.ObjectId(id).Hex()
}

func (id ObjectID) Valid() bool {
	return bson.ObjectId(id).Valid()
}

func ObjectIdHex(s string) ObjectID {
	return ObjectID(bson.ObjectIdHex(s))
}

func IsObjectIdHex(s string) bool {
	return bson.IsObjectIdHex(s)
}

func NewObjectId() ObjectID {
	return ObjectID(bson.NewObjectId())
}
