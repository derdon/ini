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
		if item.property == property {
			return true
		}
	}
	return false
}

// Returns a list of all section names of the config.
// Note: the order is not deterministic!
func (c *Config) GetSections() (sections []string) {
	sections = []string{}
	for section, _ := range *c {
		sections = append(sections, section)
	}
	return sections
}

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

func (c *Config) Get(section, property string) (value string, err error) {
	items, err := c.GetItems(section)
	if err != nil {
		return value, err
	}
	for _, item := range items {
		if (*item).property == property {
			return (*item).value, nil
		}
	}
	return value, NoPropertyError{property}
}

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

// The accepted values are: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// All other values result in an error.
func (c *Config) GetBool(section, property string) (value bool, err error) {
	f := func(s string) (interface{}, error) { return strconv.ParseBool(s) }
	v, err := c.GetFormatted(section, property, f)
	return v.(bool), err
}

// Gets the value of the given property in the given section and returns it to
// an integer. If the value cannot be converted to an int, an error is returned.
// For other possible error return values, see the documentation of the Get
// method.
func (c *Config) GetInt(section, property string) (value int, err error) {
	f := func(s string) (interface{}, error) { return strconv.Atoi(s) }
	v, err := c.GetFormatted(section, property, f)
	return v.(int), err
}

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

func (c *Config) GetFloat64(section, property string) (value float64, err error) {
	f := func(s string) (interface{}, error) { return strconv.ParseFloat(s, 64) }
	v, err := c.GetFormatted(section, property, f)
	return v.(float64), err
}