package main

import "fmt"

import "github.com/derdon/ini"

func main() {
	// create a new empty config
	conf := ini.NewConfig()
	// add a section
	conf.AddSection("my little section")
	// and another
	conf.AddSection("temporary section")
	// ... and remove it
	conf.RemoveSection("temporary section")
	// create a new property and initialise it with a value
	conf.Set("my little section", "temp property", "temp value")
	// ... and remove it
	conf.RemoveProperty("my little section", "temp property")
	// config now only contains the section "my little section" and no items
	fmt.Println(conf)
}
