package structs

/*
Struct for encoding json response for RENEWABLES_CURRENT and RENEWABLES_HISTORY endpoints.
 */
type CountryOutput struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Year       string  `json:"year,omitempty"` //  suppress field if not defined, such as when returning mean percentage value.
	Percentage float64 `json:"percentage"`
}

/*
Countries as stored in country cache and for interactions with restcountires API.
 */
type Country struct {
	Name    string   `json:"name"`
	IsoCode string   `json:"isoCode"`
	Borders []string `json:"borders"`
}

/*
Struct for encoding JSON response for deleting and viewing a webhook/all webhooks in Notification endpoint.
 */
type Webhook struct {
	WebhookId string `json:"webhook_id"`
	Url       string `json:"url,omitempty"`
	Country   string `json:"country,omitempty"`
	Calls     int    `json:"calls,omitempty"`
	Year      int    `json:"year,omitempty"`
}

/*
Struct for status endpoint response
*/
type Status struct {
	CountriesApi   string  `json:"countries_api"`
	NotificationDb string  `json:"notification_db"`
	Webhooks       int     `json:"webhooks"`
	Version        string  `json:"version"`
	Uptime         float64 `json:"uptime"`
}
