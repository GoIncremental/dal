package dal

import "errors"

var ErrInvalidConnectionString = errors.New("Invalid Connection String")
var ErrAlreadyConnected = errors.New("Already Connected")
var ErrNotConnected = errors.New("Not Connected")

//Query is a simple mapping type that allows
//clients to filter the data they wish to retrieve
type Query map[string]interface{}

//ID is required to identify the item being operated on
type ID interface {
	IsValid() bool
}

//Item is an inteface that defines the necessary methods
//to enable Connection to perform its jobs
type Item interface {
	GetID() ID
}

//Connection is the main interface definition for DAL
//it defines the contract
type Connection interface {
	Connect(string) (Connection, error)
	Clone() (Connection, error)
	Close() error
	Tag(string) (Connection, error)
	Create(Item) (interface{}, error)
	Read(ID) (interface{}, error)
	// Update(ID, Item) error
	// Delete(ID) error
	// Find(Query) (interface{}, error)
}
