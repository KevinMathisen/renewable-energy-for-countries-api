package handlers

import (
	"assignment2/utils/structs"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Tests the sorting of countries by ISO code
*/
func TestSortByIsoCode(t *testing.T) {

	// Create input
	input := [][]structs.CountryOutput{
		{
			{Name: "Norway", IsoCode: "NOR", Year: "2020", Percentage: 69},
			{Name: "Norway", IsoCode: "NOR", Year: "2021", Percentage: 70},
		},
		{
			{Name: "Sweden", IsoCode: "SWE", Year: "2020", Percentage: 50},
			{Name: "Sweden", IsoCode: "SWE", Year: "2021", Percentage: 51},
		},
		{
			{Name: "Finland", IsoCode: "FIN", Year: "2020", Percentage: 58},
			{Name: "Finland", IsoCode: "FIN", Year: "2021", Percentage: 59},
		},
	}

	// Create expected output
	expected := []structs.CountryOutput{
		{Name: "Finland", IsoCode: "FIN", Year: "2020", Percentage: 58},
		{Name: "Finland", IsoCode: "FIN", Year: "2021", Percentage: 59},
		{Name: "Norway", IsoCode: "NOR", Year: "2020", Percentage: 69},
		{Name: "Norway", IsoCode: "NOR", Year: "2021", Percentage: 70},
		{Name: "Sweden", IsoCode: "SWE", Year: "2020", Percentage: 50},
		{Name: "Sweden", IsoCode: "SWE", Year: "2021", Percentage: 51},
	}

	// Try to run the fuction
	output := sortByIsoCode(input)

	// check if we got the expected number of structs back
	assert.Equal(t, len(output), 6, "Wrong amount of structs returned: ", len(output))

	// Check if we got the expected output
	assert.Equal(t, output, expected, "Output is wrong")

}
