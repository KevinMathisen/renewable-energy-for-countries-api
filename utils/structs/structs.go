package structs

/*
* Struct for encoding json response for RENEWABLES_CURRENT and RENEWABLES_HISTORY endpoints.
 */
type CountryOutput struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Year       string  `json:"year,omitempty"` //  suppress field if not defined, such as when returning mean percentage value.
	Percentage float64 `json:"percentage"`
}

/*
* Countries as stored in country cache and for interactions with restcountires API.
 */
type Country struct {
	Name    string `json:"name"`
	IsoCode string `json:"isoCode"`
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

/*
* Struct for wrapping errors for standardized error handling.
*
* OrigErr: Original error message
* StatusCode: Status code to show user
* UsrMessage: Error message to show user.
* DevMessage: Error message to display in logs.
 */
type WrappedError struct {
	OrigErr    error
	StatusCode int
	UsrMessage string
	DevMessage string
}
