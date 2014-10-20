package main

func main() {

	InitLogger()
	go ProcessImg("IMG_1335.JPG");
	StartConsuming();

}
