package main

import (
	"archive/zip"
	"time"
	"fmt"
	"os"
	"io"
)

type ZIP_ERROR struct {
	Message string
}

type Jobber func()



func (err *ZIP_ERROR) Error() string{
	return err.Message
}

func ProcessZip(dataDir,loc, userName string)(string, error){
	r,err := zip.OpenReader(loc)
	if nil != err{
		LogOnError(err, "problem with zip " + loc)
		return "",err
	}
	defer r.Close()


	t:= time.Now()

	uzipDir := dataDir + "/" + userName + "/" + t.Format("2006-01-02")

	_,err = os.Stat(uzipDir)

	if nil != err{
		err = os.MkdirAll(uzipDir, os.ModePerm)
		if nil != err{
			fmt.Errorf("failed to make dir %s ", err.Error)
			return "",err
		}
	}


	fmt.Println(" unzip dir " + uzipDir)

	// copy images over
	for _,f:= range r.File{

			rc, err := f.Open()
			if err != nil {
				fmt.Println("failed to open file " + f.Name)
				return "",err
			}
			defer rc.Close()
			out, err := os.Create(uzipDir + "/" + f.Name)

			if err != nil {
				fmt.Println("failed to create file " + uzipDir + "/" + f.Name + " " + err.Error())
				return "",err;
			}
			defer out.Close()
			_, err = io.Copy(out, rc)
			if nil != err {
				fmt.Println("failed to copy file " + uzipDir + "/" + f.Name)
				return "",err
			}
	}

	return uzipDir,err

}
