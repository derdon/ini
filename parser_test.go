package ini

import (
	"strings"
	"testing"
)

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
	linereader := NewLineReader(strings.NewReader(""))
	line, err := linereader.ReadLine()
	assertErrorIsNil(err, t)
	expectLine("", line, t)
}

func TestReadlineNoNewline(t *testing.T) {
	linereader := NewLineReader(strings.NewReader("line without nl"))
	line, err := linereader.ReadLine()
	assertErrorIsNil(err, t)
	expectLine("line without nl", line, t)
}

func TestReadline(t *testing.T) {
	reader := strings.NewReader("first line\nsecond line")
	linereader := NewLineReader(reader)
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
	expectProperty("foo", item.property, t)
	expectValue("bar", item.value, t)
}

func TestParseItemWithWhitespace(t *testing.T) {
	line := "foo  = 	bar"
	item, err := parseItem(line)
	assertErrorIsNil(err, t)
	expectProperty("foo", item.property, t)
	expectValue("bar", item.value, t)
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
	expectProperty("foo", item.property, t)
	expectValue("bar \\= baz", item.value, t)
}

// TODO: test parseItem with quoted values!

func TestParseINIEmpty(t *testing.T) {
	linereader := NewLineReader(strings.NewReader(""))
	config, err := ParseINI(linereader)
	assertErrorIsNil(err, t)
	expectedConfig := make(Config)
	assertConfigMapsEqual(config, &expectedConfig, t)
}

func TestParseINIOneSection(t *testing.T) {
	linereader := NewLineReader(strings.NewReader("[section]"))
	config, err := ParseINI(linereader)
	assertErrorIsNil(err, t)
	section := make(map[string]string)
	expectedConfig := &Config{"section": section}
	assertConfigMapsEqual(config, expectedConfig, t)
}

func TestParseINITwoSections(t *testing.T) {
	fileContent := "[section one]\n[section two]"
	linereader := NewLineReader(strings.NewReader(fileContent))
	config, err := ParseINI(linereader)
	assertErrorIsNil(err, t)
	sectionOne := make(map[string]string)
	sectionTwo := make(map[string]string)
	expectedConfig := &Config{
		"section one": sectionOne,
		"section two": sectionTwo}
	assertConfigMapsEqual(config, expectedConfig, t)
}

func TestParseINISectionWithOneAssignment(t *testing.T) {
	filecontent := "[section]\nproperty=value"
	linereader := NewLineReader(strings.NewReader(filecontent))
	config, err := ParseINI(linereader)
	assertErrorIsNil(err, t)
	expectedConfig := &Config{"section": {"property": "value"}}
	assertConfigMapsEqual(config, expectedConfig, t)
}

func TestParseINIAssignmentBeforeSection(t *testing.T) {
	filecontent := "property=value\n[section]"
	linereader := NewLineReader(strings.NewReader(filecontent))
	_, err := ParseINI(linereader)
	assertErrorIsNotNil(err, t)
	if err != AssignmentOutsideSectionError {
		t.Errorf("expected AssignmentOutsideSectionError, got %v", err)
	}
}

func TestParseINIBrokenAssignment(t *testing.T) {
	filecontent := "[section]\nproperty value"
	linereader := NewLineReader(strings.NewReader(filecontent))
	_, err := ParseINI(linereader)
	assertErrorIsNotNil(err, t)
	if err != MissingEqualSignError {
		t.Errorf("expected MissingEqualSignError, got %v", err)
	}
}