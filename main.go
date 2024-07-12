package main

import "User-Management/http_interface"

func main() {
	app := &http_interface.Handler{}
	app.Initialize()
	app.Run()
}
