package main

import (
	"fmt"
	"strings"
)

import "github.com/derdon/ini"

func main() {
	filecontent := `[section one]
[another section]
foo = bar`
	linereader := ini.NewLineReader(strings.NewReader(filecontent))
	conf, err := ini.ParseINI(linereader)
	if err != nil {
		panic(err)
	}

	// print all sections, seperated by commas
	sections := conf.GetSections()
	for i, section := range sections {
		fmt.Printf("section #%d: %q\n", i+1, section)
	}
	fmt.Println()

	// error will be nil, because we know that the passed section exists
	items, _ := conf.GetItems("section one")
	fmt.Printf("items of \"section one\": %v\n\n", items)

	// print the items of the section "another section"
	items, _ = conf.GetItems("another section")
	fmt.Println("items of \"another section\": ")
	for _, item := range items {
		fmt.Printf("\tproperty: %q, value: %q\n", item.Property, item.Value)
	}
}
