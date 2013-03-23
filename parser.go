package ini

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

const newline = 10

var MissingEqualSignError = errors.New("missing equal sign")
var TooManyEqualSignsError = errors.New("too many equal signs")

var assignmentPattern = regexp.MustCompile("[^\\\\]=")

type LineReader struct {
	io.ByteReader
}

func NewLineReader(r io.ByteReader) *LineReader {
	return &LineReader{r}
}

// Read bytes from the given LineReader until a newline occurs. If the reader
// contains no newlines, its whole content is returned. If the reader is empty,
// i.e. contains no bytes at all, the empty string and no error is returned.
func (r *LineReader) ReadLine() (line string, err error) {
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

func ParseINI(reader *LineReader) (*Config, error) {
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
			// if the line is not a section, it must be an assignment
			// otherwise it's a syntax error
			item, err := parseItem(line)
			if err != nil {
				return &conf, err
			}
			if section != "" {
				conf.Set(section, item.property, item.value)
			} else {
				// assignment outside a section.
				// this is a syntax error
				return &conf, AssignmentOutsideSectionError
			}
		}
	}
	return &conf, nil
}
