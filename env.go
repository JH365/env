package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	_TagName    = "env"
	_TagDefault = "default"
)

var errWrongType = errors.New("[env] require a pointer points to a struct")

// Parse is the main entry for this library
func Parse(v interface{}) error {

	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Ptr {
		return errWrongType
	}

	val := ptr.Elem()
	if val.Kind() != reflect.Struct {
		return errWrongType
	}

	return doParse(val)
}

// doParse parses the env variables based on struct field tag
func doParse(refVal reflect.Value) error {

	refType := refVal.Type()

	for i := 0; i < refVal.NumField(); i++ {
		refFieldValue := refVal.Field(i)
		refFieldType := refType.Field(i)

		// Check if the filed is an unexported one
		if !refFieldValue.CanSet() {
			continue
		}

		// Check if the field is a pointer
		if refFieldValue.Kind() == reflect.Ptr {
			refFieldValue = refFieldValue.Elem()
		}

		if refFieldValue.Kind() == reflect.Struct {
			if err := doParse(refFieldValue); err != nil {
				return err
			}
			continue
		}

		envVal := getEnvByTag(refFieldType)
		if envVal == "" {
			continue
		}

		if err := setValue(refFieldValue, envVal); err != nil {
			fmt.Println(err)
			return err
		}

	}
	return nil
}

// getEnvByTag fetch env value based on field tag
// TODO: Adding support for more configurations like "required", "customized parse function", etc.
func getEnvByTag(refField reflect.StructField) string {

	envName := refField.Tag.Get(_TagName)
	if envName == "" {
		return ""
	}

	if envVal, ok := os.LookupEnv(envName); ok {
		return envVal
	}

	if envDefault, ok := refField.Tag.Lookup(_TagDefault); ok {
		return envDefault
	}

	return ""
}

// setValue converts env string to corresponding data type
// TODO: Adding supports for more data types
func setValue(refVal reflect.Value, val string) error {

	switch refVal.Kind() {
	case reflect.String:
		refVal.SetString(val)
	case reflect.Int:
		v, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		refVal.SetInt(v)
	// case reflect.Int8:
	// 	v, err := strconv.ParseInt(val, 10, 8)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Int16:
	// 	v, err := strconv.ParseInt(val, 10, 16)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Int32:
	// 	v, err := strconv.ParseInt(val, 10, 32)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Int64:
	// 	v, err := strconv.ParseInt(val, 10, 64)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Uint:
	// 	v, err := strconv.ParseInt(val, 10, 32)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Uint8:
	// 	v, err := strconv.ParseInt(val, 10, 8)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Uint16:
	// 	v, err := strconv.ParseInt(val, 10, 16)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Uint32:
	// 	v, err := strconv.ParseInt(val, 10, 32)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	// case reflect.Uint64:
	// 	v, err := strconv.ParseInt(val, 10, 64)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	refVal.SetInt(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		refVal.SetBool(v)
	default:
		return errors.New("unsupported value type")
	}

	return nil
}
