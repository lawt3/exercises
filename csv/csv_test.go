package csv

import (
	"encoding/csv"
	"os"
	"testing"
)

func TestUnmarshal(t *testing.T) {
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
