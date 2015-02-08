package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"os"
)


const(
	AWQ_CREDENTIAL_PATH = "/etc/photomap/.aws/credentials"
)



func PutInBucket(file string, remoteName string)(string,error){
	f,err := os.Open(file);

	if nil != err {
		return "",err;
	}
	defer f.Close()
	auth,err :=aws.GetAuth(CONF.getAwsAccessKey(),CONF.getAwsSecretKey())
	s3conn :=	s3.New(auth,aws.EUWest)
	bucket := s3conn.Bucket("photo-map")
	data,err := ioutil.ReadAll(f)
	if nil != err{
		return "",err;
	}

	InfoLog.Println(" adding filepath " + remoteName)
	err =  bucket.Put(remoteName, data, "image/jpeg", s3.AuthenticatedRead);
	if nil != err{
		return "",err;
	}
	return CONF.getAwsLocation() + remoteName,nil;
}
