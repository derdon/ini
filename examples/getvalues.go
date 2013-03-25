package main

import (
	"fmt"
	"strings"
)

import "github.com/derdon/ini"

func main() {
	filecontent := `[section]
goethe quote = Da steh ich nun, ich armer Tor! Und bin so klug als wie zuvor.
sense of life = 42
sqrt of two = 1.41421356237
is this a boolean = true
`
	linereader := ini.NewLineReader(strings.NewReader(filecontent))
	conf, err := ini.ParseINI(linereader)
	if err != nil {
		panic(err)
	}
	// We know that both the section and the property exist,
	// so the error value can be discarded
	value, _ := conf.Get("section", "goethe quote")
	fmt.Printf("the value of \"goethe quote\" is: %q\n", value)

	integer, _ := conf.GetInt("section", "sense of life")
	fmt.Printf("the value of \"sense of life\" is: %d\n", integer)

	float, _ := conf.GetFloat32("section", "sqrt of two")
	fmt.Printf("the value of \"sqrt of two\" is: %f\n", float)

	boolean, _ := conf.GetBool("section", "is this a boolean")
	fmt.Printf("the value of \"is this a boolean\" is: %t\n", boolean)
}
