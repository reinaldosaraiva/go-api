package main

import "github.com/reinaldosaraiva/go-api/configs"

func main() {
	config,_ := configs.LoadConfig(".")
	println(config.DBHost)
}
