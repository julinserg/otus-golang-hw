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

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, ve := range v {
		sb.WriteString(fmt.Sprintf("Field: %s, Err: %s", ve.Field, ve.Err.Error()))
	}
	return sb.String()
}

type AppError struct {
	Err error
}

func (r AppError) Error() string {
	return fmt.Sprintf("Err: %s", r.Err)
}

func checkStringLen(fieldValue *reflect.Value, validatorName string, validatorValue string) bool {
	if validatorName != "len" {
		return true
	}
	intVar, err := strconv.Atoi(validatorValue)
	_ = err
	return len(fieldValue.String()) == intVar
}

func checkStringContainsInSubSet(fieldValue *reflect.Value, validatorName string, validatorValue string) bool {
	if validatorName != "in" {
		return true
	}
	subsetList := strings.Split(validatorValue, ",")
	isConatain := false
	for _, st := range subsetList {
		if st == fieldValue.String() {
			isConatain = true
			break
		}
	}
	return isConatain
}

func checkStringMatchRegexp(fieldValue *reflect.Value, validatorName string, validatorValue string) bool {
	if validatorName != "regexp" {
		return true
	}
	r, err := regexp.Compile(validatorValue)
	if err != nil {
		fmt.Println("regexp error")
		return false
	}
	return r.MatchString(fieldValue.String())
}

func checkIntMatchMin(fieldValue *reflect.Value, validatorName string, validatorValue string) bool {
	if validatorName != "min" {
		return true
	}
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		fmt.Println("min error")
		return false
	}
	return fieldValue.Int() <= int64(minValue)
}

func checkIntMatchMax(fieldValue *reflect.Value, validatorName string, validatorValue string) bool {
	if validatorName != "max" {
		return true
	}
	maxValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		fmt.Println("error convert string to int")
		return false
	}
	return int64(maxValue) >= fieldValue.Int()
}

func checkIntContainsInSubSet(fieldValue *reflect.Value, validatorName string, validatorValue string) bool {
	if validatorName != "in" {
		return true
	}
	subsetList := strings.Split(validatorValue, ",")
	isConatain := false
	for _, st := range subsetList {
		stInt, err := strconv.Atoi(st)
		if err != nil {
			fmt.Println("error convert string to int")
			return false
		}
		if int64(stInt) == fieldValue.Int() {
			isConatain = true
			break
		}
	}
	return isConatain
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	checkInt := func(fieldValue *reflect.Value, validatorName string, validatorValue string, varName string) bool {
		if !checkIntMatchMin(fieldValue, validatorName, validatorValue) {
			validationErrors = append(validationErrors, ValidationError{Field: varName, Err: ValidateErrorNotMatchMin})
			return false
		}
		if !checkIntMatchMax(fieldValue, validatorName, validatorValue) {
			validationErrors = append(validationErrors, ValidationError{Field: varName, Err: ValidateErrorNotMatchMax})
			return false
		}
		if !checkIntContainsInSubSet(fieldValue, validatorName, validatorValue) {
			validationErrors = append(validationErrors, ValidationError{Field: varName, Err: ValidateErrorNotContainsInt})
			return false
		}
		return true
	}

	checkString := func(fieldValue *reflect.Value, validatorName string, validatorValue string, varName string) bool {
		if !checkStringLen(fieldValue, validatorName, validatorValue) {
			validationErrors = append(validationErrors, ValidationError{Field: varName, Err: ValidateErrorBadLength})
			return false
		}
		if !checkStringContainsInSubSet(fieldValue, validatorName, validatorValue) {
			validationErrors = append(validationErrors, ValidationError{Field: varName, Err: ValidateErrorNotContainsString})
			return false
		}
		if !checkStringMatchRegexp(fieldValue, validatorName, validatorValue) {
			validationErrors = append(validationErrors, ValidationError{Field: varName, Err: ValidateErrorNotMatchRegexp})
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
					fmt.Println("panic error")
					continue
				}

				switch fieldValue.Kind() {
				case reflect.Int:
					checkInt(&fieldValue, argList[0], argList[1], fieldType.Name)
				case reflect.String:
					checkString(&fieldValue, argList[0], argList[1], fieldType.Name)
				case reflect.Slice:
				L:
					for i := 0; i < fieldValue.Len(); i++ {
						sliceValue := fieldValue.Index(i)
						fmt.Println(sliceValue)
						switch sliceValue.Kind() {
						case reflect.Int:
							if !checkInt(&sliceValue, argList[0], argList[1], fieldType.Name) {
								fmt.Println("break 1")
								break L
							}
						case reflect.String:
							if !checkString(&sliceValue, argList[0], argList[1], fieldType.Name) {
								fmt.Println("break 2")
								break L
							}
						default:
							fmt.Println("Unsupported type")
							continue
						}
					}
				default:
					fmt.Println("Unsupported type")
					continue
				}
			}

		}
		return validationErrors

	}
	return &AppError{Err: errors.New("v not struct")}
}
