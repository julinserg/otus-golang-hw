package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: App{
				Version: "26.08.2022",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Version", Err: errors.New("bad length")}},
		},
		{
			in: App{
				Version: "26.08",
			},
			expectedErr: nil,
		},
		{
			in:          "struct",
			expectedErr: &AppError{Err: errors.New("v not struct")},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			actualError := Validate(tt.in)
			if tt.expectedErr != nil {
				require.EqualError(t, actualError, tt.expectedErr.Error())

				var eApp *AppError
				var eValidate ValidationErrors
				if errors.As(tt.expectedErr, &eApp) {
					require.Equal(t, errors.As(actualError, &eApp), true)
				} else {
					require.Equal(t, errors.As(actualError, &eValidate), true)
				}

			} else {
				require.Nil(t, actualError)
			}

		})
	}
}
