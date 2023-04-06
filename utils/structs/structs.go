package structs

/*
* Struct for encoding json response for RENEWABLES_CURRENT and RENEWABLES_HISTORY endpoints.
 */
type CountryOutput struct {
	Name       string `json:"name"`
	IsoCode    string `json:"isoCode"`
	Year       string `json:"year,omitempty"` //  suppress field if not defined, such as when returning mean percentage value.
	Percentage string `json:"percentage"`
}
