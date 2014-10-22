package main

func main() {

	InitLogger()
	LoadConfig()
	go ProcessPhotoDir()
	StartConsuming();

}
