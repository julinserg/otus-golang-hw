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

var ValidateErrorBadLength = errors.New("bad length")
var ValidateErrorNotContainsString = errors.New("not contains (string)")
var ValidateErrorNotMatchRegexp = errors.New("not match regexp")
var ValidateErrorNotMatchMin = errors.New("not match min")
var ValidateErrorNotMatchMax = errors.New("not match max")
var ValidateErrorNotContainsInt = errors.New("not contains (int)")

var AppErrorNotStruct = errors.New("v not struct")
var AppErrorBadValidatorSeparator = errors.New("bad validator separator")

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
		return ValidateErrorBadLength
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
		return ValidateErrorNotContainsString
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
		return ValidateErrorNotMatchRegexp
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
		return ValidateErrorNotMatchMin
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
		return ValidateErrorNotMatchMax
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
		return ValidateErrorNotContainsInt
	}
	return nil
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	checkInt := func(fieldValue *reflect.Value, validatorName string, validatorValue string, varName string) error {
		if err := checkIntMatchMin(fieldValue, validatorName, validatorValue); err != nil {
			if errors.Is(err, ValidateErrorNotMatchMin) {
				validationErrors = append(validationErrors, ValidationError{Field: varName, Err: err})
			} else {
				return err
			}
		}
		if err := checkIntMatchMax(fieldValue, validatorName, validatorValue); err != nil {
			if errors.Is(err, ValidateErrorNotMatchMax) {
				validationErrors = append(validationErrors, ValidationError{Field: varName, Err: err})
			} else {
				return err
			}
		}
		if err := checkIntContainsInSubSet(fieldValue, validatorName, validatorValue); err != nil {
			if errors.Is(err, ValidateErrorNotContainsInt) {
				validationErrors = append(validationErrors, ValidationError{Field: varName, Err: err})
			} else {
				return err
			}
		}
		return nil
	}

	checkString := func(fieldValue *reflect.Value, validatorName string, validatorValue string, varName string) error {
		if err := checkStringLen(fieldValue, validatorName, validatorValue); err != nil {
			if errors.Is(err, ValidateErrorBadLength) {
				validationErrors = append(validationErrors, ValidationError{Field: varName, Err: err})
			} else {
				return err
			}
		}
		if err := checkStringContainsInSubSet(fieldValue, validatorName, validatorValue); err != nil {
			if errors.Is(err, ValidateErrorNotContainsString) {
				validationErrors = append(validationErrors, ValidationError{Field: varName, Err: err})
			} else {
				return err
			}
		}
		if err := checkStringMatchRegexp(fieldValue, validatorName, validatorValue); err != nil {
			if errors.Is(err, ValidateErrorNotMatchRegexp) {
				validationErrors = append(validationErrors, ValidationError{Field: varName, Err: err})
			} else {
				return err
			}
		}
		return nil
	}

	var checkFields func(fieldValue *reflect.Value, validatorName string, validatorValue string, fieldName string) bool

	checkFields = func(fieldValue *reflect.Value, validatorName string, validatorValue string, fieldName string) bool {
		switch fieldValue.Kind() {
		case reflect.Int:
			checkInt(fieldValue, validatorName, validatorValue, fieldName)
		case reflect.String:
			checkString(fieldValue, validatorName, validatorValue, fieldName)
		case reflect.Slice:
			for i := 0; i < fieldValue.Len(); i++ {
				sliceValue := fieldValue.Index(i)
				checkFields(&sliceValue, validatorName, validatorValue, fieldName+" index "+strconv.Itoa(i))
			}
		default:
			return false

		}
		return true
	}

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
					return &AppError{Err: AppErrorBadValidatorSeparator}
				}
				if !checkFields(&fieldValue, argList[0], argList[1], fieldType.Name) {
					continue
				}
			}
		}
		return validationErrors
	}
	return &AppError{Err: AppErrorNotStruct}
}
