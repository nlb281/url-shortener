package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	config := config.MustLoad()

	fmt.Println(config)
	
	// TODO: init logger

	// TODO: init storage

	// TODO: init router

	// TODO: run server
}