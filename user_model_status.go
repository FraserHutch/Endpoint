package main

// ModelStatusCode - result codes from db operations.
type ModelStatusCode int

// Status codes returned by the model
const (
	ModelSuccess ModelStatusCode = iota
	ModelDBCreateFailure
	ModelDBGetFailure
	ModelDBUpdateFailure
	ModelDBDeleteFailure
)

var modelStatusText = map[ModelStatusCode]string{
	ModelSuccess:         "Success",
	ModelDBCreateFailure: "User create failure",
	ModelDBGetFailure:    "User get failure",
	ModelDBUpdateFailure: "User update failure",
	ModelDBDeleteFailure: "User delete failure",
}

// ModelStatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func ModelStatusText(code ModelStatusCode) string {
	return modelStatusText[code]
}
