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

	StructContainsString struct {
		Variable string `validate:"in:foo,bar"`
	}

	StructRegexpMatch struct {
		Variable string `validate:"regexp:\\d+"`
	}

	StructMatchMin struct {
		Variable int `validate:"min:18"`
	}

	StructMatchMax struct {
		Variable int `validate:"max:45"`
	}

	StructContainsInt struct {
		Variable int `validate:"in:11,12"`
	}

	StructInStringInt struct {
		Variable1 int    `validate:"in:11,12"`
		Variable2 string `validate:"in:foo,bar"`
	}

	StructSliceInt struct {
		VariableSlice []int `validate:"in:11,12"`
	}

	StructMultiValidator struct {
		Variable string `validate:"in:foo,bar|len:3"`
	}

	StructPrivate struct {
		variable string `validate:"len:3"`
	}

	StructBadSeparator struct {
		Variable string `validate:"len-3"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: StructContainsString{
				Variable: "qwerty",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ValidateErrorNotContainsString}},
		},
		{
			in: StructContainsString{
				Variable: "foo",
			},
			expectedErr: nil,
		},
		{
			in: StructRegexpMatch{
				Variable: "123",
			},
			expectedErr: nil,
		},
		{
			in: StructRegexpMatch{
				Variable: "yyy",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ValidateErrorNotMatchRegexp}},
		},
		{
			in: StructMatchMin{
				Variable: 12,
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ValidateErrorNotMatchMin}},
		},
		{
			in: StructMatchMin{
				Variable: 18,
			},
			expectedErr: nil,
		},
		{
			in: StructMatchMax{
				Variable: 60,
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ValidateErrorNotMatchMax}},
		},
		{
			in: StructMatchMax{
				Variable: 45,
			},
			expectedErr: nil,
		},
		{
			in: StructContainsInt{
				Variable: 78,
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ValidateErrorNotContainsInt}},
		},
		{
			in: StructContainsInt{
				Variable: 11,
			},
			expectedErr: nil,
		},
		{
			in: StructInStringInt{
				Variable1: 45,
				Variable2: "45",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Variable1", Err: ValidateErrorNotContainsInt},
				ValidationError{Field: "Variable2", Err: ValidateErrorNotContainsString}},
		},
		{
			in: StructSliceInt{
				VariableSlice: []int{77, 78},
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "VariableSlice index 0", Err: ValidateErrorNotContainsInt},
				ValidationError{Field: "VariableSlice index 1", Err: ValidateErrorNotContainsInt}},
		},
		{
			in: StructMultiValidator{
				Variable: "foo",
			},
			expectedErr: nil,
		},
		{
			in: StructMultiValidator{
				Variable: "fooo",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ValidateErrorNotContainsString},
				ValidationError{Field: "Variable", Err: ValidateErrorBadLength}},
		},
		{
			in: App{
				Version: "26.08.2022",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Version", Err: ValidateErrorBadLength}},
		},
		{
			in: App{
				Version: "26.08",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{4, 5, 6},
				Signature: []byte{7, 8, 9},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "Body",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 300,
				Body: "Body",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "Code", Err: ValidateErrorNotContainsInt}},
		},
		{
			in: User{
				ID:     "0a4ba8cd-b4a3-40ce-87bf-ad059468e00c",
				Name:   "name",
				Age:    18,
				Email:  "email@email.com",
				Role:   "admin",
				Phones: []string{"+0123456789"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "0a4ba8cd-b4a3-40ce-87bfad059468e00c",
				Name:   "name",
				Age:    51,
				Email:  "email@@email.com",
				Role:   "admin123",
				Phones: []string{"+0123456789", "+01234567890"},
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "ID", Err: ValidateErrorBadLength},
				ValidationError{Field: "Age", Err: ValidateErrorNotMatchMax},
				ValidationError{Field: "Email", Err: ValidateErrorNotMatchRegexp},
				ValidationError{Field: "Role", Err: ValidateErrorNotContainsString},
				ValidationError{Field: "Phones index 1", Err: ValidateErrorBadLength},
			},
		},
		{
			in:          "struct",
			expectedErr: &AppError{Err: AppErrorNotStruct},
		},
		{
			in: StructPrivate{
				variable: "fooo",
			},
			expectedErr: &ValidationErrors{ValidationError{Field: "variable", Err: ValidateErrorBadLength}},
		},
		{
			in: StructBadSeparator{
				Variable: "fooo",
			},
			expectedErr: &AppError{Err: AppErrorBadValidatorSeparator},
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
