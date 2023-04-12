package gateway

import (
	"assignment2/handlers"
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"encoding/json"
	"io"
	"strings"
)

// Map of countries
var RcCache = make(map[string]*structs.Country)

// Input a country name
func getCountry(name string) (*structs.Country, error) {

	// Check if country name is in map
	country, ok := RcCache[name]
	if ok { //Cache hit
		return country, nil
	} else { //Cache miss

		//Stitch together complete URL based on constants and input name
		urlParts := []string{constants.COUNTRIES_API_URL, constants.COUNTRY_NAME_SEARCH_PATH, name}
		url := strings.Join(urlParts, "")

		//Send get request to restcountries API
		res, err := handlers.HttpRequestFromUrl(url, "GET")
		if err != nil {
			return nil, err
		}

		//Extracts body of API response
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		//Decode response into an object
		var resObject interface{}
		err = json.Unmarshal(resBody, resObject)
		if err != nil {
			return nil, err
		}

		//Define new country struct, and fill it with data from response
		country := new(structs.Country)
		country.Name = resObject.([]interface{})[0].(map[string]interface{})["fifa"].(string)
		country.IsoCode = resObject.([]interface{})[0].(map[string]interface{})["name"].(map[string]interface{})["common"].(string)

		RcCache[country.Name] = country
	}
	// Return pointer to country
	return country, nil
}
