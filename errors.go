package ini

import (
	"errors"
	"fmt"
)

var AssignmentOutsideSectionError = errors.New(
	"attempted to use an assignment before a section definition")

var DuplicateSectionError = errors.New(
	"attempted to add a section which already exists")

var NoSectionError = errors.New(
	"attempted to set a property on a section which does not exist")

type NoPropertyError struct {
	property string
}

func (error NoPropertyError) Error() string {
	return fmt.Sprintf("No such property %q", error.property)
}
