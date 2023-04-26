package http_test_utils

import (
	"assignment2/utils/structs"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

/*
Gets data from the test URL and decodes into a slice of type CountryOutput, then returns this if
there are no errors
*/
func GetData(client http.Client, url string) ([]structs.CountryOutput, error) {
	log.Println("Testing URL: \"" + url + "\"...")

	//Sends Get request
	res, err := client.Get(url)
	if err != nil {
		log.Println("Get request to URL failed:")
		return nil, err
	}

	var resObject []structs.CountryOutput
	//Recieves values, and decodes into slice
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		log.Println("Error during decoding:")
		return nil, err
	}

	return resObject, nil
}

/*
Tests to see if the float check is equal to float mark to 13 decimal places
This is to avoid floating point errors that seem to appear around the 13th decimal place
*/
func TestPercentage(check, mark float64) string {
	checkInt := int(check * math.Pow(10, 12))
	markInt := int(mark * math.Pow(10, 12))

	if checkInt != markInt {
		return ("Mean percentage for country is not correct." +
			"\n\tExpected: " + strconv.FormatFloat(mark, 'g', -1, 64) +
			"\n\tRecieved: " + strconv.FormatFloat(check, 'g', -1, 64))
	}
	return ""
}

/*
Tests the length of given slice against the given int.
*/
func TestLen(s []structs.CountryOutput, expected int) string {
	if n := len(s); n != expected {
		return ("Wrong amount of objects returned." +
			"\n\tExpected: " + strconv.Itoa(expected) +
			"\n\tRecieved: " + strconv.Itoa(n))
	}
	return ""
}

/*
Compares the values of the CountryOutput recieved with the other values recieved in the param
*/
func TestValues(c structs.CountryOutput, code string, name string, year int, per float64) string {
	if c.IsoCode != code || !strings.EqualFold(c.Name, name) || c.Year != strconv.Itoa(year) || TestPercentage(c.Percentage, per) != "" {
		return ("Wrong object recieved." +
			"\n\tExpected: " + code + " - " + name + " - " + strconv.Itoa(year) + " - " + strconv.FormatFloat(per, 'g', -1, 64) +
			"\n\tRecieved: " + c.IsoCode + " - " + c.Name + " - " + c.Year + " - " + strconv.FormatFloat(c.Percentage, 'g', -1, 64))
	}
	return ""
}

/*
Tests the length of given slice, then tests the values of the first object, then the values of the last

	s:			the slice of objects to be operated on
	objects:	the amount of objects to be expected in s
	fISO:		the expected code of the fist object in s
	fName:		the expected name of the fist object in s
	fYear:		the expected year of the fist object in s
	fPer:		the expected percentage of the fist object in s
	lISO:		the expected code of the last object in s
	lName:		the expected name of the last object in s
	lYear:		the expected year of the last object in s
	lPer:		the expected percentage of the last object in s

	return:		a message explaining whats wrong if the test failed, empty string if test succeeded
*/
func TestLenLastFirst(s []structs.CountryOutput, objects int, fISO string, fName string, fYear int, fPer float64, lISO string, lName string, lYear int, lPer float64) string {
	//Checks that all instances of the country is recieved
	if err := TestLen(s, objects); err != "" {
		return err
	}

	//Checks that the data in the first recieved object is correct. Is case-insensitive on the country name.
	if err2 := TestValues(s[0], fISO, fName, fYear, fPer); err2 != "" {
		return err2
	}

	last := objects - 1
	//Checks that the data in the last recieved object is correct. Is case-insensitive on the country name.
	if err2 := TestValues(s[last], lISO, lName, lYear, lPer); err2 != "" {
		return err2
	}

	return ""
}

/*
Tests if fPer1 is larger than lPer1, and that fPer2 is larger than lPer2
Used to determine if a slice is sorted correctly for slices too large to list out manually
*/
func TestSortedPercentage(fPer1, lPer1, fPer2, lPer2 float64) string {
	fPer1Int := int(fPer1 * math.Pow(10, 12))
	lPer1Int := int(lPer1 * math.Pow(10, 12))
	fPer2Int := int(fPer2 * math.Pow(10, 12))
	lPer2Int := int(lPer2 * math.Pow(10, 12))

	if fPer1Int < lPer1Int || fPer2Int < lPer2Int {
		return ("The order of the sorted values is wrong.")
	}
	return ""
}

func TestSortedCodeList(check []structs.CountryOutput, mark []string) string {
	equal := true //Assume it is correct
	for i, v := range mark {
		if check[i].IsoCode != v {
			equal = false
		}
	}
	//Wrong order
	if !equal {
		return ("The list of neighbours is not correct." +
			"\n\tExpected: " + mark[0] + " - " + mark[1] + " - " + mark[2] + " - " + mark[3] +
			"\n\tRecieved: " + check[0].IsoCode + " - " + check[1].IsoCode + " - " + check[2].IsoCode + " - " + check[3].IsoCode)
	}
	return ""
}
