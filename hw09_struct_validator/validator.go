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
	ErrAppNotStruct                = errors.New("v not struct")
	ErrAppBadValidatorSeparator    = errors.New("bad validator separator")
	ErrAppTypeNotSupported         = errors.New("type not supported")
	ErrAppValidatorTagNotSupported = errors.New("validator tag not supported")
)

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, ve := range v {
		sb.WriteString(fmt.Sprintf("Field: %s, Err: %s; ", ve.Field, ve.Err.Error()))
	}
	return sb.String()
}

type AppError struct {
	Info string
	Err  error
}

func (r AppError) Error() string {
	if len(r.Info) == 0 {
		return fmt.Sprintf("Err: %s", r.Err)
	}
	return fmt.Sprintf("Err: %s - %s", r.Info, r.Err)
}

type Validator interface {
	check(fieldValue *reflect.Value, validatorValue string) error
}

type CheckStringLen struct{}

func (c *CheckStringLen) check(fieldValue *reflect.Value, validatorValue string) error {
	intVar, err := strconv.Atoi(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if len(fieldValue.String()) != intVar {
		return ErrValidateBadLength
	}
	return nil
}

type CheckStringContainsInSubSet struct{}

func (c *CheckStringContainsInSubSet) check(fieldValue *reflect.Value, validatorValue string) error {
	subsetList := strings.Split(validatorValue, ",")
	if !containsString(subsetList, fieldValue.String()) {
		return ErrValidateNotContainsString
	}
	return nil
}

type CheckStringMatchRegexp struct{}

func (c *CheckStringMatchRegexp) check(fieldValue *reflect.Value, validatorValue string) error {
	r, err := regexp.Compile(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if !r.MatchString(fieldValue.String()) {
		return ErrValidateNotMatchRegexp
	}
	return nil
}

type CheckIntMatchMin struct{}

func (c *CheckIntMatchMin) check(fieldValue *reflect.Value, validatorValue string) error {
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if fieldValue.Int() < int64(minValue) {
		return ErrValidateNotMatchMin
	}
	return nil
}

type CheckIntMatchMax struct{}

func (c *CheckIntMatchMax) check(fieldValue *reflect.Value, validatorValue string) error {
	maxValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return &AppError{Err: err}
	}
	if fieldValue.Int() > int64(maxValue) {
		return ErrValidateNotMatchMax
	}
	return nil
}

type CheckIntContainsInSubSet struct{}

func (c *CheckIntContainsInSubSet) check(fieldValue *reflect.Value, validatorValue string) error {
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

type CheckerAndError struct {
	checker      Validator
	checkerError error
}

var validatorsStringMap = map[string]CheckerAndError{
	"len":    {&CheckStringLen{}, ErrValidateBadLength},
	"in":     {&CheckStringContainsInSubSet{}, ErrValidateNotContainsString},
	"regexp": {&CheckStringMatchRegexp{}, ErrValidateNotMatchRegexp},
}

var validatorsIntMap = map[string]CheckerAndError{
	"min": {&CheckIntMatchMin{}, ErrValidateNotMatchMin},
	"max": {&CheckIntMatchMax{}, ErrValidateNotMatchMax},
	"in":  {&CheckIntContainsInSubSet{}, ErrValidateNotContainsInt},
}

func checkType(validators map[string]CheckerAndError, validationErrors *ValidationErrors, fieldValue *reflect.Value,
	validatorName string, validatorValue string, varName string,
) error {
	validatorExist := false
	for validName, validChecker := range validators {
		if validatorName == validName {
			validatorExist = true
			if err := validChecker.checker.check(fieldValue, validatorValue); err != nil {
				if errors.Is(err, validChecker.checkerError) {
					*validationErrors = append(*validationErrors, ValidationError{Field: varName, Err: err})
				} else {
					return err
				}
			}
		}
	}
	if !validatorExist {
		return &AppError{Err: ErrAppValidatorTagNotSupported, Info: validatorName}
	}
	return nil
}

func checkFields(validationErrors *ValidationErrors, fieldValue *reflect.Value,
	validatorName string, validatorValue string, fieldName string,
) error {
	typeIsSupported := false
	if fieldValue.Kind() == reflect.Int {
		typeIsSupported = true
		err := checkType(validatorsIntMap, validationErrors, fieldValue, validatorName, validatorValue, fieldName)
		if err != nil {
			return err
		}
	}
	if fieldValue.Kind() == reflect.String {
		typeIsSupported = true
		err := checkType(validatorsStringMap, validationErrors, fieldValue, validatorName, validatorValue, fieldName)
		if err != nil {
			return err
		}
	}
	if fieldValue.Kind() == reflect.Slice {
		typeIsSupported = true
		for i := 0; i < fieldValue.Len(); i++ {
			sliceValue := fieldValue.Index(i)
			err := checkFields(validationErrors, &sliceValue, validatorName, validatorValue, fieldName+" index "+strconv.Itoa(i))
			if err != nil {
				return err
			}
		}
	}
	if !typeIsSupported {
		return &AppError{Err: ErrAppTypeNotSupported, Info: fieldValue.Kind().String()}
	}
	return nil
}

func Validate(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Struct {
		return &AppError{Err: ErrAppNotStruct}
	}

	var validationErrors ValidationErrors
	structValue := reflect.ValueOf(v)
	structType := reflect.TypeOf(v)
	for i := 0; i < structValue.NumField(); i++ {
		fieldValue := structValue.Field(i)
		fieldType := structType.Field(i)
		// check only public fields
		if !fieldValue.CanInterface() {
			continue
		}
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
