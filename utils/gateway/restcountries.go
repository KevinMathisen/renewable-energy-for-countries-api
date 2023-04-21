package gateway

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"encoding/json"
	"io"
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

	//Send get request to restcountries API
	res, err := HttpRequestFromUrl(url, http.MethodGet)
	if err != nil {
		return nil, structs.NewError(err, http.StatusBadGateway, constants.DEFAULT500, "Restcountries API did not respond to request.")
	}

	//Extracts body of API response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, structs.NewError(err, http.StatusInternalServerError, constants.DEFAULT500, "Could not extract body from restcountries HTTP response.")
	}

	//Decode response into an object
	var resObject interface{}
	err = json.Unmarshal(resBody, resObject)
	if err != nil {
		return nil, structs.NewError(err, http.StatusInternalServerError, constants.DEFAULT500, "Could not decode restcountries json response.")
	}

	//Define new country struct, and fill it with data from response
	country = new(structs.Country)
	country.Name = resObject.([]interface{})[0].(map[string]interface{})[constants.USED_COUNTRY_CODE].(string)
	country.IsoCode = resObject.([]interface{})[0].(map[string]interface{})["name"].(map[string]interface{})["common"].(string)
	country.Borders = resObject.([]interface{})[0].(map[string]interface{})["borders"].([]string)

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

	//Send get request to restcountries API
	res, err := HttpRequestFromUrl(url, http.MethodGet)
	if err != nil {
		return nil, structs.NewError(err, http.StatusBadGateway, constants.DEFAULT500, "Restcountries API did not respond to request.")
	}

	//Extracts body of API response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, structs.NewError(err, http.StatusInternalServerError, constants.DEFAULT500, "Could not extract body from restcountries HTTP response.")
	}

	//Decode response into an object
	var resObject interface{}
	err = json.Unmarshal(resBody, resObject)
	if err != nil {
		return nil, structs.NewError(err, http.StatusInternalServerError, constants.DEFAULT500, "Could not decode restcountries json response.")
	}

	//Define new country struct, and fill it with data from response
	country := new(structs.Country)
	country.Name = resObject.([]interface{})[0].(map[string]interface{})[constants.USED_COUNTRY_CODE].(string)
	country.IsoCode = resObject.([]interface{})[0].(map[string]interface{})["name"].(map[string]interface{})["common"].(string)
	country.Borders = resObject.([]interface{})[0].(map[string]interface{})["borders"].([]string)

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
Get name from countries ISO code
*/
func GetNameFromIsoCode(isoCode string) (string, error) {

	country, err := GetCountryByIso(isoCode)
	if err != nil {
		return "", err
	}

	return country.Name, nil
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
