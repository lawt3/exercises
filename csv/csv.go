package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Unmarshal maps the data in a slice of slice of strings into a slice of structs.
// The first row is assumed to be the header with the column names
func Unmarshal(data [][]string, s any) error {
	sliceValPtr := reflect.ValueOf(s)
	if sliceValPtr.Kind() != reflect.Pointer {
		return errors.New("must be a pointer to a slice of structs")
	}
	sliceVal := sliceValPtr.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return errors.New("must be a pointer to a slice of structs")
	}
	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return errors.New("must be a pointer to a slice of structs")
	}

	// Assume the first row is a header
	header := data[0]
	dataCol := make(map[string]int, len(header))
	for i, name := range header {
		dataCol[name] = i
	}

	for _, record := range data[1:] {
		newVal := reflect.New(structType).Elem()
		err := unmarshalOne(record, dataCol, newVal)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, newVal))
	}

	return nil
}

func unmarshalOne(record []string, dataCol map[string]int, vv reflect.Value) error {
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		typeField := vt.Field(i)
		col, ok := dataCol[typeField.Tag.Get("csv")]
		if !ok {
			continue
		}
		val := record[col]

		field := vv.Field(i)
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(i)
		case reflect.String:
			field.SetString(val)
		case reflect.Bool:
			b, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			field.SetBool(b)
		default:
			return fmt.Errorf("cannot handle field of kind %v", field.Kind())
		}
	}

	return nil
}
