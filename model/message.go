package model

import "gopkg.in/mgo.v2/bson"
type Message struct {
	File   string
	User   string
	ResKey string
	Name   string
	FileId bson.ObjectId
}
