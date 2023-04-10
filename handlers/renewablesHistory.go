package handlers

import (
	"net/http"
)

func RenewablesHistory(w http.ResponseWriter, r *http.Request) {

}

/*
getHistoryPercentageOfRenewablesForCountries(countries, startYear, endYear, wantMean)	//History

	if startYear or endYear not specified
		set them as default values


	if countires specified and not/no mean:
		getRenewablesForCountriesByYears(isoList, startYear, endYear)					Global1

	if countires specified and mean:
		getMeanRenewablesFromCountries(isoList, startYear, endYear)						History1
			for countires
				getMeanRenewablesFromCountry(isoCode, startYear, endYEar)				History1.1


	else if no countires specified and want mean
		getMeanRenewablesFromAllCountries()												History2
			for all countires
				getMeanRenewablesFromCountry(startYear, endYEar)						History2.2

	else if no countries specified and no mean:
		getRenewablesForAllCountiresByYear(startYear, endYear)							Global2



*/
