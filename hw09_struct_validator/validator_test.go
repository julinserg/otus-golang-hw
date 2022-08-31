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

	StructBadType struct {
		Variable float64 `validate:"min:3"`
	}

	StructContainsIntErrorConvert struct {
		Variable int `validate:"in:aa,bb"`
	}
)

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in: StructContainsString{
			Variable: "qwerty",
		},
		expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ErrValidateNotContainsString}},
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
		expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ErrValidateNotMatchRegexp}},
	},
	{
		in: StructMatchMin{
			Variable: 12,
		},
		expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ErrValidateNotMatchMin}},
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
		expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ErrValidateNotMatchMax}},
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
		expectedErr: &ValidationErrors{ValidationError{Field: "Variable", Err: ErrValidateNotContainsInt}},
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
		expectedErr: &ValidationErrors{
			ValidationError{Field: "Variable1", Err: ErrValidateNotContainsInt},
			ValidationError{Field: "Variable2", Err: ErrValidateNotContainsString},
		},
	},
	{
		in: StructSliceInt{
			VariableSlice: []int{77, 78},
		},
		expectedErr: &ValidationErrors{
			ValidationError{Field: "VariableSlice index 0", Err: ErrValidateNotContainsInt},
			ValidationError{Field: "VariableSlice index 1", Err: ErrValidateNotContainsInt},
		},
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
		expectedErr: &ValidationErrors{
			ValidationError{Field: "Variable", Err: ErrValidateNotContainsString},
			ValidationError{Field: "Variable", Err: ErrValidateBadLength},
		},
	},
	{
		in: App{
			Version: "26.08.2022",
		},
		expectedErr: &ValidationErrors{ValidationError{Field: "Version", Err: ErrValidateBadLength}},
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
		expectedErr: &ValidationErrors{ValidationError{Field: "Code", Err: ErrValidateNotContainsInt}},
	},
	{
		in: User{
			ID:     "0a4ba8cd-b4a3-40ce-87bf-ad059468e00c",
			Name:   "name",
			Age:    18,
			Email:  "email@email.com",
			Role:   "admin",
			Phones: []string{"+0123456789"},
			meta:   nil,
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
			meta:   nil,
		},
		expectedErr: &ValidationErrors{
			ValidationError{Field: "ID", Err: ErrValidateBadLength},
			ValidationError{Field: "Age", Err: ErrValidateNotMatchMax},
			ValidationError{Field: "Email", Err: ErrValidateNotMatchRegexp},
			ValidationError{Field: "Role", Err: ErrValidateNotContainsString},
			ValidationError{Field: "Phones index 1", Err: ErrValidateBadLength},
		},
	},
	{
		in:          "struct",
		expectedErr: &AppError{Err: ErrAppNotStruct},
	},
	{
		in: StructPrivate{
			variable: "fooo",
		},
		expectedErr: &ValidationErrors{ValidationError{Field: "variable", Err: ErrValidateBadLength}},
	},
	{
		in: StructBadSeparator{
			Variable: "fooo",
		},
		expectedErr: &AppError{Err: ErrAppBadValidatorSeparator},
	},
	{
		in: StructBadType{
			Variable: 5.0,
		},
		expectedErr: &AppError{Err: ErrAppTypeNotSupported},
	},
	{
		in: StructContainsIntErrorConvert{
			Variable: 12,
		},
		expectedErr: &AppError{Err: errors.New("strconv.Atoi: parsing \"aa\": invalid syntax")},
	},
}

func TestValidate(t *testing.T) {
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
