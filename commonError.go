package main

import "fmt"

func FailOnError(err error, msg string) {
	if err != nil {
		ErrorLog.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func LogOnError(err error, msg string){
	if err !=nil{
		ErrorLog.Println(err, msg)
	}
}
