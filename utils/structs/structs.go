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

/*
* Struct for decoding JSON POST request for Registration of Webhooks in Notification endpoint.
 */
type NewWebhook struct {
	Url     string `json:"url"`
	Country string `json:"country"`
	Calls   string `json:"calls"`
}

/*
* Struct for encoding JSON response for deleting and viewing a webhook/all webhooks in Notification endpoint.
 */
type Webhook struct {
	WebhookId string `json:"webook_id"`
	Url       string `json:"url,omitempty"`
	Country   string `json:"country,omitempty"`
	Calls     string `json:"calls,omitempty"`
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
