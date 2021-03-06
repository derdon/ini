package ini

import (
	"fmt"
	"strings"
	"testing"
)

type tests []struct {
	in  string
	out string
}

func assertIsSection(sectionName string, t *testing.T) {
	if !isSection(sectionName) {
		t.Errorf("%q is not a valid section name", sectionName)
	}
}

func assertIsNotSection(sectionName string, t *testing.T) {
	if isSection(sectionName) {
		t.Errorf("%q is a valid section name", sectionName)
	}
}

func expectLine(expectedLine, actualLine string, t *testing.T) {
	if actualLine != expectedLine {
		t.Errorf("expected line %q, got %q", actualLine, expectedLine)
	}
}

func expectProperty(expectedProperty, actualProperty string, t *testing.T) {
	if actualProperty != expectedProperty {
		errmsg := "expected property %q, got %q"
		t.Errorf(errmsg, expectedProperty, actualProperty)
	}
}

func expectValue(expectedValue, actualValue string, t *testing.T) {
	if actualValue != expectedValue {
		errmsg := "expected value %q, got %q"
		t.Errorf(errmsg, expectedValue, actualValue)
	}
}

func TestReadlineEmpty(t *testing.T) {
	linereader := newLineReader(strings.NewReader(""))
	line, err := linereader.ReadLine()
	assertErrorIsNil(err, t)
	expectLine("", line, t)
}

func TestReadlineNoNewline(t *testing.T) {
	linereader := newLineReader(strings.NewReader("line without nl"))
	line, err := linereader.ReadLine()
	assertErrorIsNil(err, t)
	expectLine("line without nl", line, t)
}

func TestReadline(t *testing.T) {
	reader := strings.NewReader("first line\nsecond line")
	linereader := newLineReader(reader)
	line, err := linereader.ReadLine()
	assertErrorIsNil(err, t)
	expectLine("first line\n", line, t)
}

func TestIsSectionEmptyString(t *testing.T) {
	assertIsNotSection("", t)
}

func TestIsSectionEmptySection(t *testing.T) {
	assertIsNotSection("[]", t)
}

func TestIsSectionOneLetterName(t *testing.T) {
	assertIsSection("[a]", t)
}

func TestIsSectionValid(t *testing.T) {
	assertIsSection("[validsection]", t)
}

func TestUnescapeControlCharactersEscaped(t *testing.T) {
	var unescapeTests = tests{
		{`a\rb`, "a\rb"},
		{`c\nd`, "c\nd"},
		{`e\tf`, "e\tf"},
		{`g\=h`, "g=h"},
		{`k\\l`, `k\l`}}
	var unescapedValue string
	for _, test := range unescapeTests {
		unescapedValue = unescapeControlCharacters(test.in)
		if unescapedValue != test.out {
			t.Errorf("expected %q, got %q", test.out, unescapedValue)
		}
	}
}

func TestUnescapeControlCharactersUnsupportedEscape(t *testing.T) {
	unescapedValue := unescapeControlCharacters(`\#`)
	if unescapedValue != `\#` {
		t.Errorf(`expected \#, got %q`, unescapedValue)
	}
}

func TestParseItemEmptyString(t *testing.T) {
	line := ""
	_, err := parseItem(line)
	assertErrorIsNotNil(err, t)
	if err != MissingEqualSignError {
		t.Errorf("expected MissingEqualSignError, got %v", err)
	}
}

func TestParseItemSimpleValid(t *testing.T) {
	line := "foo=bar"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.Property, t)
	expectValue("bar", item.Value, t)
}

func TestParseItemWithWhitespace(t *testing.T) {
	line := "foo  = 	bar"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.Property, t)
	expectValue("bar", item.Value, t)
}

func TestParseItemUnescapedEqualSign(t *testing.T) {
	line := "foo = bar = baz"
	_, err := parseItem(line)
	assertErrorIsNotNil(err, t)
	if err != TooManyEqualSignsError {
		t.Errorf("expected TooManyEqualSignsError, got %v", err)
	}
}

func TestParseItemWithEscapedEqualSign(t *testing.T) {
	line := "foo = bar \\= baz"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.Property, t)
	expectValue("bar = baz", item.Value, t)
}

func TestParseItemWithTab(t *testing.T) {
	line := "foo = bar \\t baz"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.Property, t)
	expectValue("bar \t baz", item.Value, t)
}

func TestParseItemWithCarriageReturn(t *testing.T) {
	line := "foo = bar \\r baz"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.Property, t)
	expectValue("bar \r baz", item.Value, t)
}

func TestParseItemWithNewline(t *testing.T) {
	line := "foo = bar \\n baz"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.Property, t)
	expectValue("bar \n baz", item.Value, t)
}

func TestParseItemWithEscapedBackslash(t *testing.T) {
	line := "foo = bar \\\\ baz"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.Property, t)
	expectValue("bar \\ baz", item.Value, t)
}

// TODO: test parseItem with quoted values!
func TestParseItemWithDoubleQuotes(t *testing.T) {}

func TestParseINIEmpty(t *testing.T) {
	config, err := NewConfigFromString("")
	assertErrorIsNil(err, t)
	expectedConfig := make(Config)
	assertConfigMapsEqual(config, &expectedConfig, t)
}

func TestParseINIComments(t *testing.T) {
	examples := []string{
		"; a little comment passing by",
		"	  ;this one starts with some whitespace",
		"# with a hash",
		" 	 	 # whitespace plus hash"}
	expectedConfig := make(Config)
	for _, input := range examples {
		config, err := NewConfigFromString(input)
		assertErrorIsNil(err, t)
		assertConfigMapsEqual(config, &expectedConfig, t)
	}
}

func TestParseINIOneSection(t *testing.T) {
	config, err := NewConfigFromString("[section]")
	assertErrorIsNil(err, t)
	section := make(map[string]string)
	expectedConfig := &Config{"section": section}
	assertConfigMapsEqual(config, expectedConfig, t)
}

func TestParseINITwoSections(t *testing.T) {
	config, err := NewConfigFromString("[section one]\n[section two]")
	assertErrorIsNil(err, t)
	sectionOne := make(map[string]string)
	sectionTwo := make(map[string]string)
	expectedConfig := &Config{
		"section one": sectionOne,
		"section two": sectionTwo}
	assertConfigMapsEqual(config, expectedConfig, t)
}

func TestParseINISectionWithOneAssignment(t *testing.T) {
	config, err := NewConfigFromString("[section]\nproperty=value")
	assertErrorIsNil(err, t)
	expectedConfig := &Config{"section": {"property": "value"}}
	assertConfigMapsEqual(config, expectedConfig, t)
}

func TestParseINIAssignmentBeforeSection(t *testing.T) {
	_, err := NewConfigFromString("property=value\n[section]")
	assertErrorIsNotNil(err, t)
	if err != AssignmentOutsideSectionError {
		t.Errorf("expected AssignmentOutsideSectionError, got %v", err)
	}
}

func TestParseINIBrokenAssignment(t *testing.T) {
	_, err := NewConfigFromString("[section]\nproperty value")
	assertErrorIsNotNil(err, t)
	if err != MissingEqualSignError {
		t.Errorf("expected MissingEqualSignError, got %v", err)
	}
}

func TestConfigStringEmpty(t *testing.T) {
	stringedConfig := NewConfig().String()
	if expectedStr := ""; stringedConfig != expectedStr {
		t.Errorf("expected %q, got %q", expectedStr, stringedConfig)
	}
}

func TestConfigStringOneSection(t *testing.T) {
	c, err := NewConfigFromString("[section]")
	assertErrorIsNil(err, t)
	stringedConfig := c.String()
	if expectedStr := "[section]"; stringedConfig != expectedStr {
		t.Errorf("expected %q, got %q", expectedStr, stringedConfig)
	}
}

func TestStringSectionWithItem(t *testing.T) {
	c, err := NewConfigFromString("[section]\nfoo	=bar")
	assertErrorIsNil(err, t)
	expectedStr := "[section]\nfoo = bar"
	if stringedConfig := c.String(); stringedConfig != expectedStr {
		t.Errorf("expected %q, got %q", expectedStr, stringedConfig)
	}
}

func ExampleConfig_String() {
	conf, _ := NewConfigFromString("[section]\n\tfoo 	\t	= bar  	")
	fmt.Println(conf)
	// Output:
	// [section]
	// foo = bar
}
