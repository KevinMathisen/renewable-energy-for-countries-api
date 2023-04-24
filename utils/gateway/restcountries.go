package gateway

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"encoding/json"
	"net/http"
	"strings"
)

// Map of countries that link country ISO codes to their respective structs containing all information.
var rcCache = make(map[string]*structs.Country)

/*
* Returns a country struct based on the ISO code of the country.
 */
func GetCountryByIso(iso string) (*structs.Country, error) {

	// Check if country ISO is in map
	country, ok := rcCache[iso]
	if ok { //Cache hit
		return country, nil
	} //Cache miss

	//Stitch together complete URL based on constants and input name
	urlParts := []string{constants.COUNTRIES_API_URL, constants.COUNTRY_CODE_SEARCH_PATH, iso}
	url := strings.Join(urlParts, "")

	country, err := getCountry(url)
	if err != nil {
		return nil, err
	}

	rcCache[country.IsoCode] = country

	// Return pointer to country
	return country, nil
}

/*
* Function takes a country name as input. It will convert it into an iso code, and then
 */
func GetCountryByName(name string) (*structs.Country, error) {

	// Check if country name is in map
	for _, v := range rcCache {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			return v, nil
		}
	}

	//Stitch together complete URL based on constants and input name
	urlParts := []string{constants.COUNTRIES_API_URL, constants.COUNTRY_NAME_SEARCH_PATH, name}
	url := strings.Join(urlParts, "")

	country, err := getCountry(url)
	if err != nil {
		return nil, err
	}

	rcCache[country.IsoCode] = country

	// Return pointer to country
	return country, nil
}

/*
Get isocode from country name
*/
func GetIsoCodeFromName(countryName string) (string, error) {

	country, err := GetCountryByName(countryName)
	if err != nil {
		return "", err
	}

	return country.IsoCode, nil
}

/*
* Gets the neighbours of input country.
 */
func GetNeighbours(isoCode string) ([]string, error) {
	country, err := GetCountryByIso(isoCode)
	if err != nil {
		return nil, err
	}

	return country.Borders, nil
}

func getInterface(url string) ([]map[string]interface{}, error) {

	//Send get request to API
	res, err := HttpRequestFromUrl(url, http.MethodGet)
	if err != nil {
		return nil, structs.NewError(err, http.StatusBadGateway, constants.DEFAULT500, "Restcountries API did not respond to request.")
	}

	//Decode response into an object
	var resObject []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resObject)
	if err != nil {
		return nil, structs.NewError(err, http.StatusInternalServerError, constants.DEFAULT500, "Could not decode restcountries json response.")
	}

	return resObject, nil
}

func getCountry(url string) (*structs.Country, error) {

	// Get response from API
	resObject, err := getInterface(url)
	if err != nil {
		return nil, err
	}

	//Define new country struct, and fill it with data from response
	country := new(structs.Country)
	country.IsoCode = resObject[0][constants.USED_COUNTRY_CODE].(string)
	country.Name = resObject[0]["name"].(map[string]interface{})["common"].(string)
	country.Borders = getCountryBorder(resObject)

	return country, nil
}

/*
Get a list of all the country ISO codes that border the country given
*/
func getCountryBorder(resObject []map[string]interface{}) []string {
	var borders []string
	// For each border, save border as a string to the list
	for _, border := range resObject[0]["borders"].([]interface{}) {
		borders = append(borders, border.(string))
	}
	return borders
}
