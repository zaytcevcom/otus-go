package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStr := ""
	for _, e := range v {
		errStr += fmt.Sprintf("Field: %s Error: %s\n", e.Field, e.Err.Error())
	}
	return errStr
}

type TagValidationError struct {
	Err error
}

func (e TagValidationError) Error() string {
	return e.Err.Error()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return errors.New("value is not struct")
	}

	valErrors := make(ValidationErrors, 0)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTag := val.Type().Field(i).Tag.Get("validate")
		if fieldTag == "" {
			continue
		}

		vErr := validateWithTag(val.Type().Field(i).Name, fieldVal, fieldTag)

		if vErr != nil {
			var tagErrors ValidationErrors
			if errors.As(vErr, &tagErrors) {
				valErrors = append(valErrors, tagErrors...)
			} else {
				return vErr
			}
		}
	}

	if len(valErrors) > 0 {
		return valErrors
	}

	return nil
}

func validateWithTag(fieldName string, field reflect.Value, tag string) error {
	valErrors := make(ValidationErrors, 0)
	tags := strings.Split(tag, "|")

	for _, t := range tags {
		data := strings.Split(t, ":")

		var err error
		var vErrors []ValidationError

		switch {
		case field.Kind() == reflect.Int:
			vErrors, err = validateInt(fieldName, field, data)

		case field.Kind() == reflect.String:
			vErrors, err = validateString(fieldName, field, data)

		case field.Kind() == reflect.Slice:
			vErrors, err = validateSlice(fieldName, field, tag)

		default:
			continue
		}

		if err != nil {
			return err
		}
		valErrors = append(valErrors, vErrors...)
	}

	if len(valErrors) > 0 {
		return valErrors
	}

	return nil
}

func validateInt(fieldName string, field reflect.Value, data []string) ([]ValidationError, error) {
	valErrors := make([]ValidationError, 0)
	num := field.Int()
	switch data[0] {
	case "min":
		min, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			return valErrors, err
		}
		if num < min {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("number less than minimum %v", min),
			})
		}
	case "max":
		max, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			return valErrors, err
		}
		if num > max {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("number greater than maximum %v", max),
			})
		}
	case "in":
		options := strings.Split(data[1], ",")
		in := false
		valueErr := false
		for _, option := range options {
			value, err := strconv.ParseInt(option, 10, 64)
			if err != nil {
				valueErr = true
				break
			}
			if num == value {
				in = true
				break
			}
		}
		if valueErr {
			return valErrors, nil
		}
		if !in {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("number not in set %v", options),
			})
		}
	}
	return valErrors, nil
}

func validateString(fieldName string, field reflect.Value, data []string) ([]ValidationError, error) {
	valErrors := make([]ValidationError, 0)
	str := field.String()
	switch data[0] {
	case "len":
		expectedLen, err := strconv.Atoi(data[1])
		if err != nil {
			return valErrors, err
		}
		if len(str) != expectedLen {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("string length not equal to %v", expectedLen),
			})
		}
	case "regexp":
		r, err := regexp.Compile(data[1])
		if err != nil {
			return valErrors, err
		}
		if !r.MatchString(str) {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("string does not match regexp %v", data[1]),
			})
		}
	case "in":
		options := strings.Split(data[1], ",")
		in := false
		for _, option := range options {
			if str == option {
				in = true
				break
			}
		}
		if !in {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("string not in set %v", options),
			})
		}
	}
	return valErrors, nil
}

func validateSlice(fieldName string, field reflect.Value, tag string) ([]ValidationError, error) {
	valErrors := make([]ValidationError, 0)
	for i := 0; i < field.Len(); i++ {
		if err := validateWithTag(fieldName, field.Index(i), tag); err != nil {
			return valErrors, err
		}
	}
	return valErrors, nil
}
