package main

import "fmt"

import "github.com/derdon/ini"

func main() {
	filecontent := "[my section]\nsome property = with a value"
	conf, err := ini.NewConfigFromString(filecontent)
	if err != nil {
		panic(fmt.Sprintf("Error: could not parse ini file. %s", err))
	}
	fmt.Println(conf)
}
