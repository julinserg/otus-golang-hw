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

var (
	ErrValidateBadLength         = errors.New("bad length")
	ErrValidateNotContainsString = errors.New("not contains (string)")
	ErrValidateNotMatchRegexp    = errors.New("not match regexp")
	ErrValidateNotMatchMin       = errors.New("not match min")
	ErrValidateNotMatchMax       = errors.New("not match max")
	ErrValidateNotContainsInt    = errors.New("not contains (int)")
)

var (
	ErrAppNotStruct             = errors.New("v not struct")
	ErrAppBadValidatorSeparator = errors.New("bad validator separator")
	ErrAppTypeNotSupported      = errors.New("type not supported")
)

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, ve := range v {
		sb.WriteString(fmt.Sprintf("Field: %s, Err: %s; ", ve.Field, ve.Err.Error()))
	}
	return sb.String()
}

type AppError struct {
	Err error
}

func (r AppError) Error() string {
	return fmt.Sprintf("Err: %s", r.Err)
}

func checkStringLen(fieldValue *reflect.Value, validatorName string, validatorValue string) error {
	if validatorName != "len" {
		return nil
	}
	intVar, err := strconv.Atoi(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if len(fieldValue.String()) != intVar {
		return ErrValidateBadLength
	}
	return nil
}

func checkStringContainsInSubSet(fieldValue *reflect.Value, validatorName string, validatorValue string) error {
	if validatorName != "in" {
		return nil
	}
	subsetList := strings.Split(validatorValue, ",")
	isConatain := false
	for _, st := range subsetList {
		if st == fieldValue.String() {
			isConatain = true
			break
		}
	}
	if !isConatain {
		return ErrValidateNotContainsString
	}
	return nil
}

func checkStringMatchRegexp(fieldValue *reflect.Value, validatorName string, validatorValue string) error {
	if validatorName != "regexp" {
		return nil
	}
	r, err := regexp.Compile(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if !r.MatchString(fieldValue.String()) {
		return ErrValidateNotMatchRegexp
	}
	return nil
}

func checkIntMatchMin(fieldValue *reflect.Value, validatorName string, validatorValue string) error {
	if validatorName != "min" {
		return nil
	}
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if !(fieldValue.Int() >= int64(minValue)) {
		return ErrValidateNotMatchMin
	}
	return nil
}

func checkIntMatchMax(fieldValue *reflect.Value, validatorName string, validatorValue string) error {
	if validatorName != "max" {
		return nil
	}
	maxValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if !(fieldValue.Int() <= int64(maxValue)) {
		return ErrValidateNotMatchMax
	}
	return nil
}

func checkIntContainsInSubSet(fieldValue *reflect.Value, validatorName string, validatorValue string) error {
	if validatorName != "in" {
		return nil
	}
	subsetList := strings.Split(validatorValue, ",")
	isConatain := false
	for _, st := range subsetList {
		stInt, err := strconv.Atoi(st)
		if err != nil {
			return &AppError{Err: err}
		}
		if int64(stInt) == fieldValue.Int() {
			isConatain = true
			break
		}
	}
	if !isConatain {
		return ErrValidateNotContainsInt
	}
	return nil
}

func checkInt(validationErrors *ValidationErrors, fieldValue *reflect.Value,
	validatorName string, validatorValue string, varName string,
) error {
	if err := checkIntMatchMin(fieldValue, validatorName, validatorValue); err != nil {
		if errors.Is(err, ErrValidateNotMatchMin) {
			*validationErrors = append(*validationErrors, ValidationError{Field: varName, Err: err})
		} else {
			return err
		}
	}
	if err := checkIntMatchMax(fieldValue, validatorName, validatorValue); err != nil {
		if errors.Is(err, ErrValidateNotMatchMax) {
			*validationErrors = append(*validationErrors, ValidationError{Field: varName, Err: err})
		} else {
			return err
		}
	}
	if err := checkIntContainsInSubSet(fieldValue, validatorName, validatorValue); err != nil {
		if errors.Is(err, ErrValidateNotContainsInt) {
			*validationErrors = append(*validationErrors, ValidationError{Field: varName, Err: err})
		} else {
			return err
		}
	}
	return nil
}

func checkString(validationErrors *ValidationErrors, fieldValue *reflect.Value,
	validatorName string, validatorValue string, varName string,
) error {
	if err := checkStringLen(fieldValue, validatorName, validatorValue); err != nil {
		if errors.Is(err, ErrValidateBadLength) {
			*validationErrors = append(*validationErrors, ValidationError{Field: varName, Err: err})
		} else {
			return err
		}
	}
	if err := checkStringContainsInSubSet(fieldValue, validatorName, validatorValue); err != nil {
		if errors.Is(err, ErrValidateNotContainsString) {
			*validationErrors = append(*validationErrors, ValidationError{Field: varName, Err: err})
		} else {
			return err
		}
	}
	if err := checkStringMatchRegexp(fieldValue, validatorName, validatorValue); err != nil {
		if errors.Is(err, ErrValidateNotMatchRegexp) {
			*validationErrors = append(*validationErrors, ValidationError{Field: varName, Err: err})
		} else {
			return err
		}
	}
	return nil
}

func checkFields(validationErrors *ValidationErrors, fieldValue *reflect.Value,
	validatorName string, validatorValue string, fieldName string,
) error {

	switch fieldValue.Kind() {
	case reflect.Int:
		err := checkInt(validationErrors, fieldValue, validatorName, validatorValue, fieldName)
		if err != nil {
			return err
		}
	case reflect.String:
		err := checkString(validationErrors, fieldValue, validatorName, validatorValue, fieldName)
		if err != nil {
			return err
		}
	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
			sliceValue := fieldValue.Index(i)
			err := checkFields(validationErrors, &sliceValue, validatorName, validatorValue, fieldName+" index "+strconv.Itoa(i))
			if err != nil {
				return err
			}
		}
	default:
		return &AppError{Err: ErrAppTypeNotSupported}
	}
	return nil
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	if reflect.ValueOf(v).Kind() == reflect.Struct {
		structValue := reflect.ValueOf(v)
		structType := reflect.TypeOf(v)
		for i := 0; i < structValue.NumField(); i++ {
			fieldValue := structValue.Field(i)
			fieldType := structType.Field(i)
			tag := fieldType.Tag
			validate, ok := tag.Lookup("validate")
			if !ok {
				continue
			}
			paramsList := strings.Split(validate, "|")
			for _, param := range paramsList {
				argList := strings.Split(param, ":")
				if len(argList) < 2 {
					return &AppError{Err: ErrAppBadValidatorSeparator}
				}
				err := checkFields(&validationErrors, &fieldValue, argList[0], argList[1], fieldType.Name)
				if err != nil {
					return err
				}
			}
		}
		return validationErrors
	}
	return &AppError{Err: ErrAppNotStruct}
}
