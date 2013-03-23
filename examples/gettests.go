package main

import (
	"fmt"
	"strings"
)

import "github.com/derdon/ini"

func main() {
	filecontent := "[names]\nname = alice"
	linereader := ini.NewLineReader(strings.NewReader(filecontent))
	conf, err := ini.ParseINI(linereader)
	if err != nil {
		panic(err)
	}

	sectionExists := conf.HasSection("names")
	fmt.Printf("Does the section 'names' exist? %t\n", sectionExists)

	propertyNameExists := conf.HasProperty("names", "name")
	fmt.Printf("Does the property 'name' exist in the section 'names'? %t\n", propertyNameExists)

	propertyAliceExists := conf.HasProperty("names", "alice")
	fmt.Printf("Does the property 'alice' exist in the section 'names'? %t\n", propertyAliceExists)
}
