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

	Id bson.ObjectId `json:"_id" bson:"_id,omitempty"`
Name      string
Path      string
Thumb     string
LonLat    []float64
Time      time.Time
TimeStamp int64
User      string
Year      string
Tags 	  []string
Complete  bool
Img string
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
	var err error;
	fmt.Println("saving pic " + pic.Id)
	if "" != pic.Id {
		_,err = c.UpsertId(pic.Id,pic)
	}else {
		err = c.Insert(pic)
	}
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
	err := c.Find(bson.M{"img": name,"user":user}).One(&result)
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

func (pic * Picture)GetPictureByIdAndUser(id, user string)(error, Picture){
	session := getDBSession()
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	var result Picture
	err := c.Find(bson.M{"user":user,"_id":bson.ObjectIdHex(id)}).One(&result)
	return err,result
}

func (pic * Picture)GetPicturesMissingData(user string)(error,[]Picture){
	session := getDBSession();
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	var result []Picture
	err :=c.Find(bson.M{"complete":false}).All(&result)
	return err,result
}

func (pic * Picture) DeletePictureById(id, user string)(error){
	session := getDBSession()
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	err := c.Remove(bson.M{"user":user,"_id":bson.ObjectIdHex(id)});
	return err;
}

func (pic * Picture) GetPicturesByLonLatAndUser(lon,lat float64, dist int, user string)(error,[]Picture){
	//as distance is in meters we want multiply the sent value
	dist = dist * 1000;
	session := getDBSession()
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	var results []Picture
	err := c.Find(bson.M{
		"lonlat": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lon, lat},
				},
				"$maxDistance": dist,
			},
		},
		"user":user,
	}).All(&results)
	return err,results
}

func (pic * Picture) GetPicturesByTagAndUser(tag , user string)(error,[]Picture){
	session := getDBSession()
	defer session.Close()
	c := session.DB(conf.CONF.GetDbName()).C(PIC_COLLECTION)
	var results []Picture
	err := c.Find(bson.M{"user":user,"tags":bson.M{"$in":[]string{tag}}}).All(&results)
	return err, results
}
