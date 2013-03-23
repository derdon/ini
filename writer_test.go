package ini

import "testing"

func TestAddNewSection(t *testing.T) {
	conf := make(Config)
	section := "section"
	if conf.HasSection(section) {
		t.Errorf("%#v has a section called %q", section)
	}
	err := conf.AddSection("section")
	assertErrorIsNil(err, t)
	if !conf.HasSection(section) {
		t.Errorf("%#v still has no section called %q", section)
	}
}

func TestAddExistingSection(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	err := conf.AddSection("section")
	assertErrorIsNotNil(err, t)
	if err != DuplicateSectionError {
		t.Errorf("expected DuplicateSectionError, got %v", err)
	}
}

func TestRemoveNonExistingSection(t *testing.T) {
	conf := make(Config)
	err := conf.RemoveSection("doesnotexist")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestRemoveExistingSection(t *testing.T) {
	conf := &Config{
		"section":  {"prop": "val"},
		"section2": make(map[string]string)}
	err := conf.RemoveSection("section")
	assertErrorIsNil(err, t)
	expectedConf := &Config{"section2": make(map[string]string)}
	assertConfigMapsEqual(conf, expectedConf, t)
}

func TestRemovePropertyFromNonExistingSection(t *testing.T) {
	conf := make(Config)
	err := conf.RemoveProperty("section", "prop")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestRemoveNonExistingProperty(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	err := conf.RemoveProperty("section", "doesnotexist")
	assertErrorIsNotNil(err, t)
	expectedError := NoPropertyError{"doesnotexist"}
	if err != expectedError {
		t.Errorf("expected %#v, got %#v", expectedError, err)
	}
}

func TestRemoveExistingProperty(t *testing.T) {
	conf := &Config{"section": {"prop": "value"}}
	err := conf.RemoveProperty("section", "prop")
	assertErrorIsNil(err, t)
	expectedConf := &Config{"section": make(map[string]string)}
	assertConfigMapsEqual(conf, expectedConf, t)
}

func TestSetPropertyValueMissingSection(t *testing.T) {
	conf := make(Config)
	err := conf.Set("section", "property", "value")
	assertErrorIsNotNil(err, t)
	if err != NoSectionError {
		t.Errorf("expected NoSectionError, got %v", err)
	}
}

func TestSetPropertyValueExistingSection(t *testing.T) {
	conf := &Config{"section": make(map[string]string)}
	err := conf.Set("section", "property", "value")
	assertErrorIsNil(err, t)
	expectedConf := &Config{"section": {"property": "value"}}
	assertConfigMapsEqual(conf, expectedConf, t)
}
