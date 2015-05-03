package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"time"
	"github.com/maleck13/photoProcessor/conf"
	"github.com/maleck13/photoProcessor/logger"
)

type Picture struct {
Name      string
Path      string
Thumb     string
LonLat    []float64
Time      time.Time
TimeStamp int64
User      string
Year      string
}



const (
	DB_NAME        = "photomap"
	PIC_COLLECTION = "pictures"
)

func getDBSession() *mgo.Session {
	fmt.Println("get db session")
	session, err := mgo.Dial(conf.CONF.GetMongoHost())
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session

}

func (pic *Picture) Save() error {
	session := getDBSession()
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	err := c.Insert(pic)
	if err != nil {
		logger.ErrorLog.Println(err.Error())
		return err
	}
	return nil
}

func (pic * Picture)FindByNameAndUser(name, user string) (error, Picture) {
	session := getDBSession()
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	result := Picture{}
	err := c.Find(bson.M{"name": name,"user":user}).One(&result)
	return err, result
}

func (pic * Picture)GetPictureDateRange(user string)(error , []string){
	session := getDBSession()
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	var result [] string
	err := c.Find(bson.M{"user":user}).Distinct("year",&result)
	return err,result
}

func (pic * Picture)GetPicturesInRange(user string , from, to int64)(error,[]Picture){
	session := getDBSession();
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	var result []Picture
	fmt.Print(bson.M{"user":user,"timestamp":bson.M{"$gte": to, "$lte": from}})
	err :=c.Find(bson.M{"user":user,"timestamp":bson.M{"$gte": to, "$lte": from}}).All(&result)
	return err,result
}