package main

import (
	"fmt"
	"strings"
)

import "github.com/derdon/ini"

func main() {
	filecontent := "[my section]\nsome property = with a value"
	reader := strings.NewReader(filecontent)
	linereader := ini.NewLineReader(reader)
	conf, err := ini.ParseINI(linereader)
	if err != nil {
		errmsg := fmt.Sprintf("Error: could not parse ini file. %s", err)
		panic(errmsg)
	}
	fmt.Printf("%#v\n", *conf)
}
