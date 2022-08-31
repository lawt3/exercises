package csv

import (
	"encoding/csv"
	"os"
	"reflect"
	"strings"
	"testing"
)

// TODO: remove
func TestRead(t *testing.T) {
	f, err := os.Open("testdata/data.csv")
	if err != nil {
		t.Fatal("Unable to read data file", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Fatal("Unable to close data file", err)
		}
	}(f)

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatal("Unable to read records", err)
	}

	t.Log(records)
}

func TestUsage(t *testing.T) {
	// TODO: Example of using this package
}

// TODO: how would you refactor this?
func TestUnmarshal(t *testing.T) {
	data := `
name,age,dob,city,has_pet
Harry,"50",23/08/1965,Melbourne,true
"Robert ""Bob"" Thompson",69,31/12/1999,Sydney,false
Tom,2,28/02/1993,Brisbane,true
`
	type Schema struct {
		Name   string `csv:"name"`
		Age    uint8  `csv:"age"`
		Dob    string `csv:"dob"`
		City   string `csv:"city"`
		HasPet bool   `csv:"has_pet"`
	}

	expected := []Schema{
		{"Harry", 50, "23/08/1965", "Melbourne", true},
		{"Robert \"Bob\" Thompson", 69, "31/12/1999", "Sydney", false},
		{"Tom", 2, "28/02/1993", "Brisbane", true},
	}

	r := csv.NewReader(strings.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatal("Failed to read records", err)
	}
	if len(records) != 4 {
		t.Fatal("Read wrong number of records", len(records))
	}

	var s []Schema
	err = Unmarshal(records, &s)
	if err != nil {
		t.Fatal("Failed to unmarshal records", records, err)
	}
	if len(s) != 4 {
		t.Fatal("Unmarshalled the wrong number of records", len(s))
	}

	for i, record := range s {
		if !reflect.DeepEqual(record, expected[i]) {
			t.Fatal("Unmarshalled struct didn't match expectations", s, expected[i])
		}
	}
}
