package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Custom struct {
		Value string `validate:"regexp:^\\d+$|len:20"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          UserRole("test"),
			expectedErr: errors.New("value is not struct"),
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Konstantin",
				Age:    27,
				Email:  "zaydisk@ya.ru",
				Role:   UserRole("admin"),
				Phones: []string{"79999999999", "78888888888"},
				meta:   json.RawMessage(`{"meta":"meta"}`),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "12345678901234567890123456789012345",
				Name:   "Konstantin",
				Age:    57,
				Email:  "zaydisk@ya.ru",
				Role:   UserRole("admin"),
				Phones: []string{"79999999999", "78888888888"},
				meta:   json.RawMessage(`{"meta":"meta"}`),
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: errors.New("string length not equal to 36")},
				ValidationError{Field: "Age", Err: errors.New("number greater than maximum 50")},
			},
		},
		{
			in: User{
				ID:     "1",
				Name:   "Konstantin",
				Age:    17,
				Email:  "zaydisk",
				Role:   UserRole("user"),
				Phones: []string{"900", "90349"},
				meta:   json.RawMessage(`{"meta":"meta"}`),
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: errors.New("string length not equal to 36")},
				ValidationError{Field: "Age", Err: errors.New("number less than minimum 18")},
				ValidationError{Field: "Email", Err: errors.New("string does not match regexp ^\\w+@\\w+\\.\\w+$")},
				ValidationError{Field: "Role", Err: errors.New("string not in set [admin stuff]")},
				ValidationError{Field: "Phones", Err: errors.New("string length not equal to 11")},
			},
		},
		{
			in: App{Version: "1.2.20"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: errors.New("string length not equal to 5")},
			},
		},
		{
			in:          App{Version: "1.2.1"},
			expectedErr: nil,
		},
		{
			in:          Token{Header: []byte("Header"), Payload: []byte("Payload"), Signature: []byte("Signature")},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 404, Body: "Not Found"},
			expectedErr: nil,
		},
		{
			in: Response{Code: 403, Body: "Forbidden"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: errors.New("number not in set [200 404 500]")},
			},
		},
		{
			in:          Custom{Value: "12345678901234567890"},
			expectedErr: nil,
		},
		{
			in: Custom{Value: "123456789012345678901234567890"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Value", Err: errors.New("string length not equal to 20")},
			},
		},
		{
			in: Custom{Value: "a2345678901234567890"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Value", Err: errors.New("string does not match regexp ^\\d+$")},
			},
		},
		{
			in: Custom{Value: "a23456789012345678901"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Value", Err: errors.New("string does not match regexp ^\\d+$")},
				ValidationError{Field: "Value", Err: errors.New("string length not equal to 20")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
			}

			if err == nil {
				return
			}

			var valErrors ValidationErrors
			var expValErrors ValidationErrors
			if errors.As(err, &valErrors) && errors.As(tt.expectedErr, &expValErrors) {
				if len(valErrors) != len(expValErrors) {
					t.Errorf("unexpected number of validation errors: expected %v, got %v\n", len(expValErrors), len(valErrors))
					return
				}

				for i, ve := range valErrors {
					eve := expValErrors[i]
					if ve.Field != eve.Field || ve.Err.Error() != eve.Err.Error() {
						t.Errorf("expected error: %+v, got: %+v", eve, ve)
					}
				}
			} else if err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error: %+v, got: %+v", tt.expectedErr, err)
			}

			_ = tt
		})
	}
}
