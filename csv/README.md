# CSV Marshaller

> *One of the best parts of developing with Go is being able to take advantage of its standard library. Like Python, it has a "batteries" included philosophy, providing many of the tools that you need to build an application.* — Jon Bodner

Go ships with the `encoding/csv` package which reads and writes CSV files. However, it doesn't provide any facilities for converting this data to or from structs. This package implements a marshaller for CSV data based on the `encoding/json` package. 

Under the hood, it uses `reflection` to inspect struct tags.

## Background
The Go standard library has support for *marshalling* and *unmarshalling* JSON through the `encoding/JSON` package. It works by using `reflection` to identify struct tags which inform how that struct should be converted to/from JSON.

For example, if we wanted to read and write JSON objects of the following form:
```JSON
{
  "id":"12345",
  "date_ordered":"2020-05-01T13:01:02Z",
  "customer_id":"3",
  "items":[{"id":"xyz123","name":"Thing 1"},{"id":"abc789","name":"Thing 2"}]
}
```

We would define types like:
```Go
type Order struct {
    ID string               `json:"id"`
    DateOrdered time.TIME   `json:"date_ordered"`
    CustomerID string       `json:"customer_id"`
    Items []Item            `json:"items"`
}

type Item struct {
    ID string   `json:"id"`
    Name string `json:"name"`
}
```

The `Unmarshal` and `Marshal` functions then turn slices of bytes into structs and vice versa.
```Go
var o Order
err := json.Unmarshal([]byte(data), &o)
if err != nil {
	return nil
}
// ...
out, err := json.Marshal(o)
```

You could just pass a `map[string]interface{}` to `json.Marshal` and `json.Unmarshal` but this loses on the strong typing Go offers. It helps to document the expected data and the types of the expected data.

## Usage
This package is follows the design of `encoding/json` and provides two functions: `csv.Marshal` and `csv.Unmarshal`.