package ini

// Add a new section to the config. If a section with this name already exists,
// the error DuplicateSectionError is returned and the section won't be added.
func (c *Config) AddSection(section string) error {
	if c.HasSection(section) {
		return DuplicateSectionError
	}
	(*c)[section] = make(map[string]string)
	return nil
}

// Remove the given section from the config. If the section does not exist, the
// error NoSectionError is returned.
func (c *Config) RemoveSection(section string) error {
	if !c.HasSection(section) {
		return NoSectionError
	}
	delete(*c, section)
	return nil
}

// Remove the given property from the passed section. If noch such section
// exists, NoSectionError will be returned. If the section exists but not the
// property, NoPropertyError will be returned.
func (c *Config) RemoveProperty(section, property string) error {
	items, err := c.GetItems(section)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.Property == property {
			delete((*c)[section], item.Property)
			return nil
		}
	}
	return NoPropertyError{property}
}

// Set the given property in the given section to the passed value. Attempting
// to set values in non-existing sections will return NoSectionError. If the
// property does not exist yet, it will be created; otherwise its value will be
// overwritten.
func (c *Config) Set(section, property, value string) error {
	if !c.HasSection(section) {
		return NoSectionError
	}
	(*c)[section][property] = value
	return nil
}
