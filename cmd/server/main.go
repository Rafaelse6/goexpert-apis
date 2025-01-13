package main

import (
	"github.com/Rafaelse6/goexpert/9-APIS/configs"
)

func main() {
	config, _ := configs.LoadConfig(".0")
	println(config.DBDriver)
}
