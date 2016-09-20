package main

import (
	"fmt"
)

func main() {
	var app TeleRDD
	fmt.Println(NO_DATE + " App launched")

	// Initialize app
	if err := app.Init(); err != nil {
		fmt.Println(err.Error())
	}
	defer app.CloseDB()

	// Launch server
	if err := app.Start(); err != nil {
		fmt.Println(err.Error())
	}
}
