package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

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

func checkLenString(validationErrors *ValidationErrors, fieldValue *reflect.Value, fieldType *reflect.StructField, validatorName string, validatorValue string) {
	if validatorName != "len" {
		return
	}
	intVar, err := strconv.Atoi(validatorValue)
	_ = err
	if len(fieldValue.String()) != intVar {
		fmt.Println("error length string ", fieldValue.String(), fieldType.Name)
		*validationErrors = append(*validationErrors, ValidationError{Field: fieldType.Name, Err: errors.New("bad length")})
	}
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
					fmt.Println("panic error")
					continue
				}
				switch fieldValue.Kind() {
				case reflect.Int:
				case reflect.String:
					checkLenString(&validationErrors, &fieldValue, &fieldType, argList[0], argList[1])

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
