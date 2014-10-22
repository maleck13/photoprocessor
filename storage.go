package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
)

type Picture struct {
	Name string
	Path string
	Thumb string
	LonLat []float64
	Time time.Time
	TimeStamp int64

}

const (
	DB_NAME = "photomap"
	PIC_COLLECTION = "pictures"
)


func getDBSession() *mgo.Session{
	fmt.Println("get db session")
	session, err := mgo.Dial(CONF.GetMongoHost())
	if err != nil {
		panic(err)
	}


	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session;

}




func (pic * Picture) Save() error{
	session := getDBSession();
	defer session.Close()
	c := session.DB(DB_NAME).C(PIC_COLLECTION)
	err := c.Insert(pic)
	if err != nil {
		ErrorLog.Println(err.Error())
		return err;
	}
	return nil;
}


func FindByName(name string) (error,Picture){
	session := getDBSession();
	defer session.Close()
	c := session.DB(DB_NAME).C(PIC_COLLECTION)
	result := Picture{}
	err := c.Find(bson.M{"name": name}).One(&result)
	return err,result
}
