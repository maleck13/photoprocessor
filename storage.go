package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Picture struct {
	Name string
	Path string
	Thumb string
	LonLat []float64

}

const (
	DB_NAME = "photomap"
	PIC_COLLECTION = "pictures"
)


func getDBSession() *mgo.Session{
	fmt.Println("get db session")
	session, err := mgo.Dial("192.168.59.103")
	if err != nil {
		panic(err)
	}


	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session;

//	c := session.DB("test").C("pictures")
//	err = c.Insert(&Picture{"test", "/test/test.jpg","test/test/test.png",[]float64{52.3434,-7.3434}})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	result := Picture{}
//	err = c.Find(bson.M{"name": "test"}).One(&result)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Name:", result.Name)
}




func SavePic(pic Picture) error{
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
