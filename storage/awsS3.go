package storage

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"os"
	"github.com/maleck13/photoProcessor/logger"
	"github.com/maleck13/photoProcessor/conf"
)

const (
	AWQ_CREDENTIAL_PATH = "/etc/photomap/.aws/credentials"
)

func PutInBucket(file string, remoteName string) (string, error) {
	f, err := os.Open(file)

	if nil != err {
		return "", err
	}
	defer f.Close()
	auth, err := aws.GetAuth(conf.CONF.GetAwsAccessKey(), conf.CONF.GetAwsSecretKey())
	s3conn := s3.New(auth, aws.EUWest)
	bucket := s3conn.Bucket("photo-map")
	data, err := ioutil.ReadAll(f)
	if nil != err {
		return "", err
	}

	logger.InfoLog.Println(" adding filepath " + remoteName)
	err = bucket.Put(remoteName, data, "image/jpeg", s3.AuthenticatedRead)
	if nil != err {
		return "", err
	}
	return conf.CONF.GetAwsLocation() + remoteName, nil
}
