package main

func main() {

	InitLogger()
	LoadConfig()
	go ProcessPhotoDir(CONF.GetPhotoDir(),CONF.GetDefaultUser())
 	StartConsuming();

}
