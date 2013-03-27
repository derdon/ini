package main

import (
	"fmt"
	"os"
)

import "github.com/derdon/ini"

func main() {
	filename := "example.ini"
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Error: could not read %s. %s", filename, err))
	}
	defer file.Close()
	conf, err := ini.NewConfigFromFile(file)
	if err != nil {
		panic(fmt.Sprintf("Error: could not parse ini file. %s", err))
	}
	fmt.Printf("%#v\n", *conf)
}
