package ini

import "strconv"

type propertyConverter func(string) (interface{}, error)

// Returns true if the config contains a section with the given name, otherwise
// false.
func (c *Config) HasSection(section string) bool {
	_, exists := (*c)[section]
	return exists
}

// Returns true if a) the given section exists and b) the given property can be
// found within the section. Otherwise false is returned.
func (c *Config) HasProperty(section, property string) bool {
	items, err := c.GetItems(section)
	if err != nil {
		return false
	}
	for _, item := range items {
		if item.Property == property {
			return true
		}
	}
	return false
}

// Returns a list of all section names of the config.
// Note: the order is not deterministic!
func (c *Config) GetSections() (sections []string) {
	sections = []string{}
	for section := range *c {
		sections = append(sections, section)
	}
	return sections
}

// Get a slice of *Item structs from the given section. The order of the
// elements in the returned slice is not determined. If the section does not
// exist, NoSectionError is returned.
func (c *Config) GetItems(section string) (items []*Item, err error) {
	items = []*Item{}
	if !c.HasSection(section) {
		return items, NoSectionError
	}
	for property, value := range (*c)[section] {
		items = append(items, &Item{property, value})
	}
	return items, nil
}

// Get the value of the passed property in the given section. If the section
// does not exist, NoSectionError is returned. If the property does not exist
// in the given section, NoPropertyError is returned.
func (c *Config) Get(section, property string) (value string, err error) {
	items, err := c.GetItems(section)
	if err != nil {
		return value, err
	}
	for _, item := range items {
		if (*item).Property == property {
			return (*item).Value, nil
		}
	}
	return value, NoPropertyError{property}
}

// Get the value of the passed property in the given section. If either the
// section doesn't exist or the property doesn't exist in the given section,
// the given default value `default_` is returned.
func (c *Config) GetDefault(section, property, default_ string) (value string) {
	value, err := c.Get(section, property)
	if err != nil {
		value = default_
	}
	return value
}

// Get the value of the passed property in the given section, apply the given
// function f to it and return the function's return values. The function must
// have the signature ``func(s string) (value interface{}, err error)``. If the
// passed section does not exist, the error NoSectionError will be returned.
// If the property does not exist within this section, NoPropertyError will be
// returned. If there was a different error returned, it came from the passed
// function.
func (c *Config) GetFormatted(section, property string, f propertyConverter) (value interface{}, err error) {
	strValue, err := c.Get(section, property)
	if err != nil {
		emptyValue, _ := f("")
		return emptyValue, err
	}
	value, err = f(strValue)
	if err != nil {
		emptyValue, _ := f("")
		return emptyValue, err
	}
	return value, nil
}

// Gets the value of the given property in the given section and returns it as
// a boolean value.
// The accepted values are: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// All other values result in an error.
func (c *Config) GetBool(section, property string) (value bool, err error) {
	f := func(s string) (interface{}, error) {
		return strconv.ParseBool(s)
	}
	v, err := c.GetFormatted(section, property, f)
	return v.(bool), err
}

// Gets the value of the given property in the given section and returns it as
// an integer. If the value cannot be converted to an int, an error is returned.
// For other possible error return values, see the documentation of the Get
// method.
func (c *Config) GetInt(section, property string) (value int, err error) {
	f := func(s string) (interface{}, error) {
		return strconv.Atoi(s)
	}
	v, err := c.GetFormatted(section, property, f)
	return v.(int), err
}

// Gets the value of the given property in the given section and returns it as
// a float32.
func (c *Config) GetFloat32(section, property string) (value float32, err error) {
	f := func(s string) (interface{}, error) {
		value, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return float32(0.0), err
		}
		return float32(value), nil
	}
	v, err := c.GetFormatted(section, property, f)
	return v.(float32), err
}

// Gets the value of the given property in the given section and returns it as
// a float64.
func (c *Config) GetFloat64(section, property string) (value float64, err error) {
	f := func(s string) (interface{}, error) {
		return strconv.ParseFloat(s, 64)
	}
	v, err := c.GetFormatted(section, property, f)
	return v.(float64), err
}
