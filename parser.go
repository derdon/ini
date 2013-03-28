package ini

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const newline = 10

var MissingEqualSignError = errors.New("missing equal sign")
var TooManyEqualSignsError = errors.New("too many equal signs")

var assignmentPattern = regexp.MustCompile("[^\\\\]=")

type lineReader struct {
	io.ByteReader
}

// create a new LineReader struct from any given io.ByteReader
func newLineReader(r io.ByteReader) *lineReader {
	return &lineReader{r}
}

// Read bytes from the given LineReader until a newline occurs. If the reader
// contains no newlines, its whole content is returned. If the reader is empty,
// i.e. contains no bytes at all, the empty string and no error is returned.
func (r *lineReader) ReadLine() (line string, err error) {
	var bytes []byte
	var b byte
	for b != newline {
		b, err = r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return line, err
		}
		bytes = append(bytes, b)
	}
	return string(bytes), nil
}

// A section is a string that start with an open bracket [, ends with an open
// bracket ] and has at least one character between those brackets.
func isSection(line string) bool {
	return (strings.HasPrefix(line, "[") &&
		strings.HasSuffix(line, "]") &&
		len(line) > 2)
}

type Item struct {
	Property string
	Value    string
}

// An assignment is of the form `name=value`. Whitespace before and after the
// equal sign is ignored. Equals signs within the value must be quoted or
// escaped with the backslash.
func parseItem(line string) (item *Item, err error) {
	matches := assignmentPattern.FindAllStringIndex(line, -1)
	if matches == nil {
		return item, MissingEqualSignError
	}
	if len(matches) > 1 {
		return item, TooManyEqualSignsError
	}
	loc := matches[0]
	property := strings.TrimSpace(line[:loc[0]+1])
	value := strings.TrimSpace(line[loc[1]:])
	item = &Item{property, value}
	return
}

// a Config map maps from section names to maps of assignments
type Config map[string]map[string]string

// Get a new empty config. This is equivalent to:
// 	NewConfigFromString("")
func NewConfig() *Config {
	c := make(Config)
	return &c
}

// Create a new *Config from a string. This is a shortcut for:
//	NewConfigFromByteReader(strings.NewReader(s))
func NewConfigFromString(s string) (*Config, error) {
	return NewConfigFromByteReader(strings.NewReader(s))
}

// Create a new *Config from a file. This is a shortcut for:
//	NewConfigFromByteReader(bufio.NewReader(file))
func NewConfigFromFile(file *os.File) (*Config, error) {
	return NewConfigFromByteReader(bufio.NewReader(file))
}

// Create a new *Config from a ByteReader.
func NewConfigFromByteReader(reader io.ByteReader) (*Config, error) {
	return parseINI(newLineReader(reader))
}

// Parse the given *LineReader to a *Config. If the reader is empty, an empty
// *Config and no error will be returned. Errors may occur when one assignment
// does not belong to any section, i.e. if it was written before the first
// section was declared. Other errors are syntax errors: Examples for syntax
// errors are: no equals sign in an assignment, more than one unescaped equal
// sign in an assignment.
func parseINI(reader *lineReader) (*Config, error) {
	conf := make(Config)
	var line string
	var err error
	var section string
	for {
		line, err = reader.ReadLine()
		if err != nil {
			return &conf, err
		}
		if line == "" {
			// stop reading at EOF
			break
		}
		trimmedLine := strings.TrimSpace(line)
		if isSection(trimmedLine) {
			section = strings.Trim(trimmedLine, "[]")
			conf.AddSection(section)
		} else {
			// If the line is not a section, it must be an
			// assignment. Otherwise it's a syntax error
			item, err := parseItem(line)
			if err != nil {
				return &conf, err
			}
			if section != "" {
				conf.Set(section, item.Property, item.Value)
			} else {
				// assignment outside a section.
				// this is a syntax error
				return &conf, AssignmentOutsideSectionError
			}
		}
	}
	return &conf, nil
}

// Return a normalized representation of the config. The order of sections and
// the order of assignments within each section is non-deterministic. Each
// section declaration begins with an open bracket [ and end with a closing
// bracket ] plus a newline \n. An Assignment starts with a property, followed
// by an equal sign which is enclosed in spaces and ends with a value and a
// newline.
func (c *Config) String() string {
	buf := new(bytes.Buffer)
	for _, section := range c.GetSections() {
		buf.WriteString(fmt.Sprintf("[%s]\n", section))
		// error can be ignored because the section surely exists
		items, _ := c.GetItems(section)
		for _, item := range items {
			buf.WriteString(fmt.Sprintf("%s = %s\n", item.Property, item.Value))
		}
	}
	// TODO: this looks inefficient and ugly. find some better way to cut of
	// the last byte of buf
	// remove trailing linebreak to be consistent with empty *Config values
	return strings.TrimSpace(string(buf.Bytes()))
}
