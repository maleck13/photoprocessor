package main

func main() {

	InitLogger()
	LoadConfig()
	go ProcessImg("IMG_1424.JPG");
	StartConsuming();

}
