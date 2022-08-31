package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Unmarshal maps the data in a slice of slice of strings into a slice of structs.
// The first row is assumed to be the header with the column names.
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

// Marshal maps a slice of structs to a slice of slice of strings.
// The first row written is the header with the column names.
func Marshal(s any) ([][]string, error) {
	sliceVal := reflect.ValueOf(s)
	if sliceVal.Kind() != reflect.Slice {
		return nil, errors.New("must be a slice of structs")
	}
	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return nil, errors.New("must be a slice of structs")
	}

	var out [][]string
	header := marshalHeader(structType)
	out = append(out, header)

	for i := 0; i < sliceVal.Len(); i++ {
		record, err := marshalOne(sliceVal.Index(i))
		if err != nil {
			return nil, err
		}
		out = append(out, record)
	}

	return out, nil
}

func marshalHeader(vt reflect.Type) []string {
	var header []string
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		if tag, ok := field.Tag.Lookup("csv"); ok {
			header = append(header, tag)
		}
	}

	return header
}

func marshalOne(vv reflect.Value) ([]string, error) {
	var record []string
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		if _, ok := vt.Field(i).Tag.Lookup("csv"); !ok {
			continue
		}

		fieldVal := vv.Field(i)
		switch fieldVal.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			record = append(record, strconv.FormatInt(fieldVal.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			record = append(record, strconv.FormatUint(fieldVal.Uint(), 10))
		case reflect.String:
			record = append(record, fieldVal.String())
		case reflect.Bool:
			record = append(record, strconv.FormatBool(fieldVal.Bool()))
		default:
			return nil, fmt.Errorf("cannot handle field of kind %v", fieldVal.Kind())
		}
	}

	return record, nil
}
