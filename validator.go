package mimir

import (
	"errors"
	"fmt"
	"reflect"
)

// SafeNilString safe nil string
func SafeNilString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}

// SafeNilInt64 safe nil string
func SafeNilInt64(value *int64) int64 {
	if value != nil {
		return *value
	}
	return 0
}

// Validate validate struct transformed into array of interface
func Validate(attributes []interface{}, m interface{}) error {
	var errStr string
	for index, attr := range attributes {
		if attr == "" || attr == nil {
			field := reflect.Indirect(reflect.ValueOf(m)).Type().Field(index).Tag.Get("json")
			errStr += fmt.Sprintf("[%s field must has value],", field)
		}
	}
	if errStr == "" {
		return nil
	}
	return errors.New(errStr)
}

// ValidateAllowNull validate struct transformed into array of interface
func ValidateAllowNull(attributes []interface{}, m interface{}) error {
	var errStr string
	for index, attr := range attributes {
		if attr == nil {
			field := reflect.Indirect(reflect.ValueOf(m)).Type().Field(index).Tag.Get("json")
			errStr += fmt.Sprintf("[%s field must has value],", field)
			println(field)
		}
	}
	if errStr == "" {
		return nil
	}
	return errors.New(errStr)
}
