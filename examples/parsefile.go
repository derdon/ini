package main

import (
	"bufio"
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
	defer f.Close()
	reader := bufio.NewReader(file)
	linereader := ini.NewLineReader(reader)
	conf, err := ini.ParseINI(linereader)
	if err != nil {
		panic(fmt.Sprintf("Error: could not parse ini file. %s", err))
	}
	fmt.Printf("%#v\n", *conf)
}
