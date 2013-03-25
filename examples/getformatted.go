package main

import (
	"fmt"
	"strings"
)

import "github.com/derdon/ini"

func parseCommaSeperatedList(s string) (interface{}, error) {
	values := []string{}
	for _, value := range strings.Split(s, ",") {
	    values = append(values, strings.TrimSpace(value))
	}
	return values, nil
}

func main() {
	filecontent := `[section]
fruits = apples, bananas, pears`
	linereader := ini.NewLineReader(strings.NewReader(filecontent))
	conf, err := ini.ParseINI(linereader)
	if err != nil {
		panic(err)
	}
	values, err := conf.GetFormatted("section", "fruits", parseCommaSeperatedList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("values: %#v\n", values)
}
