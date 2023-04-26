package structs_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"assignment2/utils/structs"
)

func TestCreateCountryOutputFromData(t *testing.T) {

	isoCode := "NOR"
	startYear := 1965
	endYear := 2021
	//Declare test date. This was is the expected format.
	data := map[string]interface{}{
		"1965": 67.87996,
		"1966": 65.3991,
		"1967": 66.591644,
		"1968": 67.13724,
		"1969": 63.88058,
		"1970": 61.510117,
		"1971": 63.9665,
		"1972": 64.29581,
		"1973": 65.582184,
		"1974": 68.31412,
		"1975": 67.97407,
		"1976": 67.33344,
		"1977": 64.26859,
		"1978": 64.17826,
		"1979": 65.86465,
		"1980": 65.252464,
		"1981": 68.59788,
		"1982": 69.16161,
		"1983": 71.88228,
		"1984": 71.13633,
		"1985": 69.52919,
		"1986": 66.977295,
		"1987": 68.43144,
		"1988": 69.97985,
		"1989": 71.44203,
		"1990": 72.44774,
		"1991": 71.44005,
		"1992": 71.865555,
		"1993": 71.17737,
		"1994": 69.031494,
		"1995": 70.81212,
		"1996": 66.01224,
		"1997": 66.31446,
		"1998": 67.218765,
		"1999": 68.73755,
		"2000": 72.39789,
		"2001": 67.58246,
		"2002": 69.30982,
		"2003": 63.816036,
		"2004": 64.23876,
		"2005": 69.73603,
		"2006": 66.73525,
		"2007": 69.2405,
		"2008": 69.89543,
		"2009": 67.78923,
		"2010": 65.47019,
		"2011": 66.30012,
		"2012": 70.095116,
		"2013": 67.50864,
		"2014": 68.88728,
		"2015": 68.87519,
		"2016": 69.86629,
		"2017": 69.260994,
		"2018": 68.85805,
		"2019": 67.08509,
		"2020": 70.96306,
		"2021": 71.558365,
		"name": "Norway",
	}

	// Call the function being tested
	output, err := structs.CreateCountryOutputFromData(data, isoCode, startYear, endYear)
	jsonOutput, _ := json.Marshal(output)

	// Check for errors
	if err != nil {
		t.Errorf("CreateCountryOutputFromData() returned error: %v", err)
	}

	// Check that the response writer contains the expected output
	expectedOutput := `[{"name":"Norway","isoCode":"NOR","year":"1965","percentage":67.87996},{"name":"Norway","isoCode":"NOR","year":"1966","percentage":65.3991},{"name":"Norway","isoCode":"NOR","year":"1967","percentage":66.591644},{"name":"Norway","isoCode":"NOR","year":"1968","percentage":67.13724},{"name":"Norway","isoCode":"NOR","year":"1969","percentage":63.88058},{"name":"Norway","isoCode":"NOR","year":"1970","percentage":61.510117},{"name":"Norway","isoCode":"NOR","year":"1971","percentage":63.9665},{"name":"Norway","isoCode":"NOR","year":"1972","percentage":64.29581},{"name":"Norway","isoCode":"NOR","year":"1973","percentage":65.582184},{"name":"Norway","isoCode":"NOR","year":"1974","percentage":68.31412},{"name":"Norway","isoCode":"NOR","year":"1975","percentage":67.97407},{"name":"Norway","isoCode":"NOR","year":"1976","percentage":67.33344},{"name":"Norway","isoCode":"NOR","year":"1977","percentage":64.26859},{"name":"Norway","isoCode":"NOR","year":"1978","percentage":64.17826},{"name":"Norway","isoCode":"NOR","year":"1979","percentage":65.86465},{"name":"Norway","isoCode":"NOR","year":"1980","percentage":65.252464},{"name":"Norway","isoCode":"NOR","year":"1981","percentage":68.59788},{"name":"Norway","isoCode":"NOR","year":"1982","percentage":69.16161},{"name":"Norway","isoCode":"NOR","year":"1983","percentage":71.88228},{"name":"Norway","isoCode":"NOR","year":"1984","percentage":71.13633},{"name":"Norway","isoCode":"NOR","year":"1985","percentage":69.52919},{"name":"Norway","isoCode":"NOR","year":"1986","percentage":66.977295},{"name":"Norway","isoCode":"NOR","year":"1987","percentage":68.43144},{"name":"Norway","isoCode":"NOR","year":"1988","percentage":69.97985},{"name":"Norway","isoCode":"NOR","year":"1989","percentage":71.44203},{"name":"Norway","isoCode":"NOR","year":"1990","percentage":72.44774},{"name":"Norway","isoCode":"NOR","year":"1991","percentage":71.44005},{"name":"Norway","isoCode":"NOR","year":"1992","percentage":71.865555},{"name":"Norway","isoCode":"NOR","year":"1993","percentage":71.17737},{"name":"Norway","isoCode":"NOR","year":"1994","percentage":69.031494},{"name":"Norway","isoCode":"NOR","year":"1995","percentage":70.81212},{"name":"Norway","isoCode":"NOR","year":"1996","percentage":66.01224},{"name":"Norway","isoCode":"NOR","year":"1997","percentage":66.31446},{"name":"Norway","isoCode":"NOR","year":"1998","percentage":67.218765},{"name":"Norway","isoCode":"NOR","year":"1999","percentage":68.73755},{"name":"Norway","isoCode":"NOR","year":"2000","percentage":72.39789},{"name":"Norway","isoCode":"NOR","year":"2001","percentage":67.58246},{"name":"Norway","isoCode":"NOR","year":"2002","percentage":69.30982},{"name":"Norway","isoCode":"NOR","year":"2003","percentage":63.816036},{"name":"Norway","isoCode":"NOR","year":"2004","percentage":64.23876},{"name":"Norway","isoCode":"NOR","year":"2005","percentage":69.73603},{"name":"Norway","isoCode":"NOR","year":"2006","percentage":66.73525},{"name":"Norway","isoCode":"NOR","year":"2007","percentage":69.2405},{"name":"Norway","isoCode":"NOR","year":"2008","percentage":69.89543},{"name":"Norway","isoCode":"NOR","year":"2009","percentage":67.78923},{"name":"Norway","isoCode":"NOR","year":"2010","percentage":65.47019},{"name":"Norway","isoCode":"NOR","year":"2011","percentage":66.30012},{"name":"Norway","isoCode":"NOR","year":"2012","percentage":70.095116},{"name":"Norway","isoCode":"NOR","year":"2013","percentage":67.50864},{"name":"Norway","isoCode":"NOR","year":"2014","percentage":68.88728},{"name":"Norway","isoCode":"NOR","year":"2015","percentage":68.87519},{"name":"Norway","isoCode":"NOR","year":"2016","percentage":69.86629},{"name":"Norway","isoCode":"NOR","year":"2017","percentage":69.260994},{"name":"Norway","isoCode":"NOR","year":"2018","percentage":68.85805},{"name":"Norway","isoCode":"NOR","year":"2019","percentage":67.08509},{"name":"Norway","isoCode":"NOR","year":"2020","percentage":70.96306},{"name":"Norway","isoCode":"NOR","year":"2021","percentage":71.558365}]`
	if string(jsonOutput) != expectedOutput {
		t.Errorf("CreateCountryOutputFromData() returned unexpected output.\nExpected: %s\nActual: %s", expectedOutput, jsonOutput)
	}

	//Check that output is in correct dataformat
	outputType := reflect.TypeOf(output)
	expectedType := reflect.TypeOf([]structs.CountryOutput{})
	if outputType != expectedType {
		t.Errorf("CreateCountryOutputFromData() returned unexpected output. \nExpected datatype: %s\nActual datatype: %s ", expectedType, outputType)
	}

	//Check if data is outside of stard and end year bounderies
	for _, val := range output {

		intYear, err := strconv.Atoi(val.Year)
		if err != nil {
			t.Errorf("CreateCountryOutputFromData() returned unexpected output. Non-number percentage in dataset.")
		}

		if intYear > endYear || intYear < startYear {
			t.Errorf("CreateCountryOutputFromData() returned unexpected output. \nExpected Year Range: " + strconv.Itoa(endYear) + "-" + strconv.Itoa(startYear) + "\nActual year: " + strconv.Itoa(intYear))
		}
	}
}

func TestCreateMeanCountryOutputFromData(t *testing.T) {

}

func TestMean(t *testing.T) {
	//Establish a number of tests consisting of a list of values, and the correct mean value of the list.
	tests := []struct {
		input []float64
		want  float64
	}{
		{[]float64{1, 2, 3}, 2},
		{[]float64{0, 0, 0}, 0},
		{[]float64{-1, 1}, 0},
		{[]float64{}, 0},
	}

	//Feed the values into the mean function and report eventual anomalies.
	for _, test := range tests {
		got := structs.Mean(test.input)
		if got != test.want {
			t.Errorf("mean(%v) = %v; want %v", test.input, got, test.want)
		}
	}
}

func TestNewError(t *testing.T) {
	// Test case 1: Create a new error with a non-nil original error
	origErr := errors.New("original error")
	statusCode := http.StatusBadRequest
	userMsg := "Bad request"
	devMsg := "Request was not well-formed"
	err := structs.NewError(origErr, statusCode, userMsg, devMsg)
	if err == nil {
		t.Errorf("NewError() returned nil, expected non-nil error")
	}
	if err.Error() != "original error" {
		t.Errorf("NewError() returned error message '%s', expected 'original error'", err.Error())
	}
	if e, ok := err.(structs.WrappedError); ok {
		if e.StatusCode != 400 {
			t.Errorf("NewError() returned status code %d, expected %d", e.StatusCode, 400)
		}
		if e.UsrMessage != "Bad request" {
			t.Errorf("NewError() returned user message '%s', expected 'Bad request'", e.UsrMessage)
		}
		if e.DevMessage != "Request was not well-formed" {
			t.Errorf("NewError() returned dev message '%s', expected 'Request was not well-formed'", e.DevMessage)
		}
	} else {
		t.Errorf("NewError() did not return a WrappedError, expected WrappedError")
	}

	// Test case 2: Create a new error with a nil original error
	origErr = nil
	statusCode = http.StatusNotFound
	userMsg = "Not found"
	devMsg = "Requested resource could not be found"
	err = structs.NewError(origErr, statusCode, userMsg, devMsg)
	if err == nil {
		t.Errorf("NewError() returned nil, expected non-nil error")
	}
	if err.Error() != "" {
		t.Errorf("NewError() returned error message '%s', expected empty string", err.Error())
	}
	if e, ok := err.(structs.WrappedError); ok {
		if e.StatusCode != 404 {
			t.Errorf("NewError() returned status code %d, expected %d", e.StatusCode, 404)
		}
		if e.UsrMessage != "Not found" {
			t.Errorf("NewError() returned user message '%s', expected 'Not found'", e.UsrMessage)
		}
		if e.DevMessage != "Requested resource could not be found" {
			t.Errorf("NewError() returned dev message '%s', expected 'Requested resource could not be found'", e.DevMessage)
		}
	} else {
		t.Errorf("NewError() did not return a WrappedError, expected WrappedError")
	}
}

func TestMyError_Error(t *testing.T) {
	// Create a new instance of MyError with a specific error message
	errMsg := "Original Error Message"
	err := structs.NewError(errors.New(errMsg), http.StatusBadRequest, "Bad request.", "Request was not well formed")

	// Call the Error() function and compare the result to the expected error message
	if err.Error() != errMsg {
		t.Errorf("Error() returned '%s', expected '%s'", err.Error(), errMsg)
	}
}
