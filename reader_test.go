package ini

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

func TestHasSectionEmptyString(t *testing.T) {
	conf := make(Config)
	if conf.HasSection("") {
		t.Errorf("%#v has section \"\"", conf)
	}
}

func TestHasSectionValid(t *testing.T) {
	conf := &Config{"some section": make(map[string]string)}
	if !conf.HasSection("some section") {
		t.Errorf("%#v has no section %q", conf, "some section")
	}
}

func TestHasPropertyMissingSection(t *testing.T) {
	conf := make(Config)
	if conf.HasProperty("doesnotexist", "prop") {
		t.Errorf("%#v has property \"prop\"", conf)
	}
}

func TestHasPropertyNoItems(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	if conf.HasProperty("section", "prop") {
		t.Errorf("%#v has property \"prop\"", conf)
	}
}

func TestHasPropertyExists(t *testing.T) {
	conf := &Config{"section": {"prop": "val"}}
	if !conf.HasProperty("section", "prop") {
		t.Errorf("%#v has no property \"prop\"", conf)
	}
}

func TestGetSectionsEmptyConf(t *testing.T) {
	conf := make(Config)
	sections := conf.GetSections()
	if !reflect.DeepEqual(sections, []string{}) {
		t.Errorf("expected []string{}, got %#v", sections)
	}
}

func TestGetSectionsAccessible(t *testing.T) {
	conf := &Config{
		"one": make(map[string]string),
		"two": make(map[string]string)}
	sections := conf.GetSections()
	// sort the result before comparing because we don't the order when
	// receiving keys from a map
	sort.Strings(sections)
	expectedSections := []string{"one", "two"}
	if !reflect.DeepEqual(sections, expectedSections) {
		t.Errorf("expected %#v, got %#v", expectedSections, sections)
	}
}

func TestGetItemsFromNonExistingSection(t *testing.T) {
	conf := make(Config)
	_, err := conf.GetItems("section")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestGetItemsSectionWithNoItems(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	items, err := conf.GetItems("section")
	assertErrorIsNil(err, t)
	expectedItems := []*Item{}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("expected %#v, got %#v", expectedItems, items)
	}
}

func TestGetItemsAccessible(t *testing.T) {
	conf := &Config{"section": {"prop": "val"}}
	items, err := conf.GetItems("section")
	assertErrorIsNil(err, t)
	expectedItems := []*Item{&Item{"prop", "val"}}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("expected %#v, got %#v", expectedItems, items)
	}
}

func TestGetMissingSection(t *testing.T) {
	conf := make(Config)
	_, err := conf.Get("section", "property")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestGetMissingProperty(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	_, err := conf.Get("section", "property")
	assertErrorIsNotNil(err, t)
	expectedError := NoPropertyError{"property"}
	if err != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestGetExisting(t *testing.T) {
	conf := &Config{"section": {"property": "value"}}
	value, err := conf.Get("section", "property")
	assertErrorIsNil(err, t)
	expectedValue := "value"
	if value != expectedValue {
		t.Errorf("expected %q, got %q", expectedValue, value)
	}
}

func TestGetFormattedMissingSection(t *testing.T) {
	conf := make(Config)
	dummyConverter := func(s string) (interface{}, error) { return s, nil }
	_, err := conf.GetFormatted("section", "property", dummyConverter)
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestGetFormattedMissingProperty(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	dummyConverter := func(s string) (interface{}, error) { return s, nil }
	_, err := conf.GetFormatted("section", "property", dummyConverter)
	assertErrorIsNotNil(err, t)
	expectedError := NoPropertyError{"property"}
	if err != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestGetFormattedConverterReturningError(t *testing.T) {
	conf := &Config{"section": {"property": "value"}}
	customError := errors.New("my custom error")
	f := func(s string) (interface{}, error) { return "", customError }
	_, err := conf.GetFormatted("section", "property", f)
	assertErrorIsNotNil(err, t)
	if err != customError {
		t.Errorf("expected %v, got %v", customError, err)
	}
}

func TestGetFormattedValid(t *testing.T) {
	conf := &Config{"section": {"property": "value"}}
	f := func(s string) (interface{}, error) { return s + s, nil }
	value, err := conf.GetFormatted("section", "property", f)
	assertErrorIsNil(err, t)
	expectedValue := "valuevalue"
	if value != expectedValue {
		t.Errorf("expected %q, got %q", expectedValue, value)
	}
}

func TestGetBoolMissingSection(t *testing.T) {
	conf := make(Config)
	_, err := conf.GetBool("section", "property")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestGetBoolMissingProperty(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	_, err := conf.GetBool("section", "property")
	assertErrorIsNotNil(err, t)
	expectedError := NoPropertyError{"property"}
	if err != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestGetBoolInvalidValue(t *testing.T) {
	conf := &Config{"section": {"property": "that's not true!"}}
	_, err := conf.GetBool("section", "property")
	// the type of the returned error is not documented, so I guess it's
	// an implementation detail and I will therefore not rely on that
	// particular error
	assertErrorIsNotNil(err, t)
}

func TestGetBoolTrue(t *testing.T) {
	for _, value := range []string{"1", "t", "T", "TRUE", "true", "True"} {
		conf := &Config{"section": {"property": value}}
		booleanValue, err := conf.GetBool("section", "property")
		assertErrorIsNil(err, t)
		if !booleanValue {
			t.Error("expected true, got false")
		}
	}
}

func TestGetBoolFalse(t *testing.T) {
	for _, value := range []string{"0", "f", "F", "FALSE", "false", "False"} {
		conf := &Config{"section": {"property": value}}
		booleanValue, err := conf.GetBool("section", "property")
		assertErrorIsNil(err, t)
		if booleanValue {
			t.Error("expected false, got true")
		}
	}
}

func TestGetIntMissingSection(t *testing.T) {
	conf := make(Config)
	_, err := conf.GetInt("section", "property")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestGetIntMissingProperty(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	_, err := conf.GetInt("section", "property")
	assertErrorIsNotNil(err, t)
	expectedError := NoPropertyError{"property"}
	if err != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestGetIntInvalidValue(t *testing.T) {
	conf := &Config{"section": {"property": "NaN"}}
	_, err := conf.GetInt("section", "property")
	assertErrorIsNotNil(err, t)
}

func TestGetIntValid(t *testing.T) {
	conf := &Config{"section": {"property": "42"}}
	intValue, err := conf.GetInt("section", "property")
	assertErrorIsNil(err, t)
	if intValue != 42 {
		t.Errorf("expected 42, got %d", intValue)
	}
}

func TestGetFloat32MissingSection(t *testing.T) {
	conf := make(Config)
	_, err := conf.GetFloat32("section", "property")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestGetFloat32MissingProperty(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	_, err := conf.GetFloat32("section", "property")
	assertErrorIsNotNil(err, t)
	expectedError := NoPropertyError{"property"}
	if err != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestGetFloat32InvalidValue(t *testing.T) {
	conf := &Config{"section": {"property": "NaF (short for: Not a Float)"}}
	_, err := conf.GetFloat32("section", "property")
	assertErrorIsNotNil(err, t)
}

func TestGetFloat32Valid(t *testing.T) {
	conf := &Config{"section": {"e": "1.718281828"}}
	floatValue, err := conf.GetFloat32("section", "e")
	assertErrorIsNil(err, t)
	if floatValue != 1.718281828 {
		t.Errorf("expected 1.718281828, got %d", floatValue)
	}
}

func TestGetFloat64MissingSection(t *testing.T) {
	conf := make(Config)
	_, err := conf.GetFloat64("section", "property")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestGetFloat64MissingProperty(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	_, err := conf.GetFloat64("section", "property")
	assertErrorIsNotNil(err, t)
	expectedError := NoPropertyError{"property"}
	if err != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestGetFloat64InvalidValue(t *testing.T) {
	conf := &Config{
		"section": {"property": "NaFe (short for: Not a Float either"}}
	_, err := conf.GetFloat64("section", "property")
	assertErrorIsNotNil(err, t)
}

func TestGetFloat64Valid(t *testing.T) {
	conf := &Config{"section": {"e": "1.718281828"}}
	floatValue, err := conf.GetFloat64("section", "e")
	assertErrorIsNil(err, t)
	if floatValue != 1.718281828 {
		t.Errorf("expected 1.718281828, got %d", floatValue)
	}
}
