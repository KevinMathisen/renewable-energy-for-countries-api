package structs

import (
	"assignment2/utils/constants"
	"net/http"
	"sort"
	"strconv"
)

/*
Creates a slice of countryOutput structs sorted by year which can be sent as response to requests
Goes through each year for a country, filters out the ones we want, and create a struct for each year

	w			- Responsewriter for error handling
	data		- Map which contain name of country and renewable percentages for all years of data
	isoCode		- isoCode of country we are creating structs for
	startYear	- The year in which we want to start returning data from
	endYear		- The year in which we want to stop returning data from

	return		- List of countryOutput structs which can be encoded into Json and sent as reponse to requests
*/
func CreateCountryOutputFromData(data map[string]interface{}, isoCode string, startYear int, endYear int) ([]CountryOutput, error) {
	var output []CountryOutput

	// Save country name as a string
	countryName := data["name"].(string)

	// For each year of country renewables
	for year, percentage := range data {

		// Ignore name field
		if year == "name" {
			continue
		}

		// Try to convert year to an int
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return output, NewError(nil, http.StatusInternalServerError, constants.DEFAULT500, "Error when creating data, could not convert year to int")

		}

		// Ignore years outside of scope defined by user
		if yearInt < startYear || yearInt > endYear {
			continue
		}

		// Create countryoutput with year and percentage
		countryOutput := CountryOutput{
			Name:       countryName,
			IsoCode:    isoCode,
			Year:       year,
			Percentage: percentage.(float64),
		}

		// Save each countryoutput to slice
		output = append(output, countryOutput)
	}

	// Sort output based on year
	sort.Slice(output, func(i, j int) bool {
		return output[i].Year < output[j].Year
	})

	return output, nil
}

/*
Creates a slice of countryOutput structs with Mean value hich can be sendt as response to requests
Goes through each year for a country, filters out the ones we want, and caulcates the mean value for all years.
Then returnes a struct with the mean value.

	w			- Responsewriter for error handling
	data		- Map which contain name of country and renewable percentages for all years of data
	isoCode		- isoCode of country we are creating struct for
	startYear	- The year in which we want to start calculating mean from
	endYear		- The year in which we want to stop calculating mean from

	return		- CountryOutput struct with no year value and mean value as percentage, can be encoded into Json and sent as reponse to requests
*/
func CreateMeanCountryOutputFromData(data map[string]interface{}, isoCode string, startYear int, endYear int) ([]CountryOutput, error) {
	var percentages []float64

	// Save country name as a string
	countryName := data["name"].(string)

	// For each year of country renewables
	for year, percentage := range data {

		// Ignore name field
		if year == "name" {
			continue
		}

		// Try to convert year to an int
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			return nil, NewError(nil, http.StatusInternalServerError, constants.DEFAULT500, "Error when creating data, could not convert year to int")
		}

		// Ignore years outside of scope defined by user
		if yearInt < startYear || yearInt > endYear {
			continue
		}

		// Add each percentage to a list of all percentages in time range
		percentages = append(percentages, percentage.(float64))
	}

	// Create a countryOutput without year and mean value as percentage
	countryOutput := CountryOutput{
		Name:       countryName,
		IsoCode:    isoCode,
		Percentage: mean(percentages),
	}

	return []CountryOutput{countryOutput}, nil
}

/*
Create a webhook struct given a map of data and webhook ID as string
*/
func CreateWebhookFromData(data map[string]interface{}, webhookID string) Webhook {
	// Create webhook struct from data
	webhook := Webhook{
		WebhookId: webhookID,
		Url:       data["url"].(string),
		Country:   data["country"].(string),
		Calls:     int(data["calls"].(int64)),
	}

	// Include year if specified
	if data["year"].(int64) != -1 {
		webhook.Year = int(data["year"].(int64))
	}

	return webhook
}

/*
Calculate mean value of list of numbers

	input	- List of float values

	return	- Average of list
*/
func mean(input []float64) float64 {
	// If there is no input
	if len(input) == 0 {
		return 0
	}

	var sum float64

	// Add all values in input to get sum
	for _, value := range input {
		sum += value
	}

	// Return mean value of input
	return sum / float64(len(input))
}
