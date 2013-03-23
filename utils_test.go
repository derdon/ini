package ini

import (
	"reflect"
	"testing"
)

func assertErrorIsNil(err error, t *testing.T) {
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func assertErrorIsNotNil(err error, t *testing.T) {
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func assertConfigMapsEqual(firstConfig, secondConfig *Config, t *testing.T) {
	if !reflect.DeepEqual(firstConfig, secondConfig) {
		t.Errorf("%#v ≠ %#v", *firstConfig, *secondConfig)
	}
}
