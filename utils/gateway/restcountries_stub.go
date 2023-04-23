package gateway

import (
	"assignment2/utils/constants"
	"assignment2/utils/structs"
	"errors"
	"net/http"
	"strings"
)

type RestCountriesStub struct {
	Countries map[string]*structs.Country
}

func (rcs *RestCountriesStub) GetCountryByIso(iso string) (*structs.Country, error) {
	country, ok := rcs.Countries[iso]
	if ok {
		return country, nil
	}
	return nil, structs.NewError(errors.New("could not get country"), http.StatusInternalServerError, constants.DEFAULT500, "Could not get country.")
}

func (rcs *RestCountriesStub) GetCountryByName(name string) (*structs.Country, error) {
	for _, v := range rcs.Countries {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			return v, nil
		}
	}
	return nil, structs.NewError(errors.New("could not get country"), http.StatusInternalServerError, constants.DEFAULT500, "Could not get country.")
}

func (rcs *RestCountriesStub) GetIsoCodeFromName(countryName string) (string, error) {
	country, err := rcs.GetCountryByName(countryName)
	if err != nil {
		return "", err
	}
	return country.IsoCode, nil
}

func (rcs *RestCountriesStub) GetNameFromIsoCode(isoCode string) (string, error) {
	country, err := rcs.GetCountryByIso(isoCode)
	if err != nil {
		return "", err
	}
	return country.Name, nil
}
