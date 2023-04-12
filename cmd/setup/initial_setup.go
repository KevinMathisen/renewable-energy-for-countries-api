package main

import (
	"assignment2/utils/constants"
	"assignment2/utils/db"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type datapoint struct {
	entity     string `json:"Entity"`
	code       string `json:"code"`
	year       string `json:"Year"`
	renewables string `json:"renewables"`
}

func main() {

	// Set up Firestore
	db.InitializeFirestore()

	// Close down client when service is done running
	defer db.CloseFirebaseClient()

	// Get data from csv file
	data := createRenewablesDataStructForAppending()

	// Add data to firestore
	_ = db.AppendDataToFirestore(data, "renewables")

}

/*
Creates a map containing all the relevant renewables data from the given csv file
*/
func createRenewablesDataStructForAppending() map[string]map[string]interface{} {

	// Open file
	file, err := os.Open(constants.RENEWABLES_CSV_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var datapoints []datapoint

	// Read csv data
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {

		log.Fatal(err)
	}

	// Export CSV data into a list of datapoint structures
	for i, line := range lines {
		if i == 0 {
			continue
		}

		row := datapoint{
			entity:     line[0],
			code:       line[1],
			year:       line[2],
			renewables: line[3],
		}

		datapoints = append(datapoints, row)
	}

	// Initialize map we will put the imported data into
	countries := make(map[string]map[string]interface{})

	// For each line of data
	for _, datapoint := range datapoints {

		// Ignore these codes, as they are not countires
		if datapoint.code == "OWID_USS" || datapoint.code == "OWID_WRL" {
			continue
		}

		// Check if country is already in map
		_, ok := countries[datapoint.code]
		if ok {
			// If country is already in map, only add the year and percentage to the country
			num, _ := strconv.ParseFloat(datapoint.renewables, 64)
			countries[datapoint.code][datapoint.year] = num

		} else if datapoint.code != "" {

			// If the country has not been added to the map, create a map containing name and one year percentage pair
			num, _ := strconv.ParseFloat(datapoint.renewables, 64)
			country := map[string]interface{}{
				"name":         datapoint.entity,
				datapoint.year: num,
			}
			// Then add the country to the map using its isoCode as its key
			countries[datapoint.code] = country
		}
	}

	return countries
}
