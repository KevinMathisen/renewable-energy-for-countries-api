package gateway

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Map of countries that link country ISO codes to their respective structs containing all information.
type RestCountriesMock struct {
	Countries map[string]*structs.Country
}

// Returns a country struct based on the ISO code of the country.
func (rcs *RestCountriesMock) GetCountryByIso(iso string) (*structs.Country, error) {
	country, ok := rcs.Countries[iso]
	if ok {
		return country, nil
	}
	return nil, structs.NewError(errors.New("could not get country"), http.StatusInternalServerError, constants.DEFAULT500, "Could not get country.")
}

// Returns a country struct based on the name of the country.
func (rcs *RestCountriesMock) GetCountryByName(name string) (*structs.Country, error) {
	for _, v := range rcs.Countries {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			return v, nil
		}
	}
	return nil, structs.NewError(errors.New("could not get country"), http.StatusInternalServerError, constants.DEFAULT500, "Could not get country.")
}

// Returns ISO code based on country name.
func (rcs *RestCountriesMock) GetIsoCodeFromName(countryName string) (string, error) {
	country, err := rcs.GetCountryByName(countryName)
	if err != nil {
		return "", err
	}
	return country.IsoCode, nil
}

// Returns country name based on ISO code.
func (rcs *RestCountriesMock) GetNameFromIsoCode(isoCode string) (string, error) {
	country, err := rcs.GetCountryByIso(isoCode)
	if err != nil {
		return "", err
	}
	return country.Name, nil
}

// Sets the country cache to the input map.
func (rcs *RestCountriesMock) SetCountryCache(countries map[string]*structs.Country) {
	rcs.Countries = countries
}

// Sets the country cache to the input map, given a JSON file.
func (rcs *RestCountriesMock) SetCountryCacheByJSON(filepath string) (map[string]*structs.Country, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var countryMap map[string]*structs.Country
	err = json.Unmarshal(byteValue, &countryMap)
	if err != nil {
		return nil, err
	}

	rcs.Countries = countryMap

	return countryMap, nil
}
