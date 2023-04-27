# Assignment2

# Overview

This service is a REST web application in Golang that provides the client with the ability to retrieve information about developments related to renewable energy production for and across countries. It uses an existing webservice, [restcountries](http://129.241.150.113:8080) for translation isoCodes and names, as well as for looking up border countries. It also uses a firestore database for saving renewables data, as well as webhooks registered, and a cache.

The service allows for notification registration using webhooks, which are invoked based on requests to specific countries (and years if specified).  

The application is dockerized and deployed using an IaaS system called Openstack on NTNUs instance called SkyHigh. See Running the assignment/Openstack instance 

## External dependencies 

The application is dependent on the following external APIs. If any of these are down the application will inform the user.

* *REST Countries API*. Endpoint: http://129.241.150.113:8080/v3.1 (Documentation: http://129.241.150.113:8080/)
* *Firebase*. Endpoint: https://console.firebase.google.com/ (Documentation: https://firebase.google.com/docs)

Dataset used for Renewables which is hosted on firebase:
* [*Renewable Energy Dataset*](https://drive.google.com/file/d/18G470pU2NRniDfAYJ27XgHyrWOThP__p/view?usp=sharing) (Authors: Hannah Ritchie, Max Roser and Pablo Rosado (2022) - "Energy". Published online at OurWorldInData.org. Retrieved from: https://ourworldindata.org/energy

The dataset reports on percentage of renewable energy in the country's energy mix over time. 

## Third-party libraries

Used the following third-party libraries: 
* Firestore, Firebase and all libraries these depend on, for interacting with the database. 
* Testify assert, for writing tests


# Completion of requirements
All the requirements are implemented, including all advanced tasks. 

* All endpoints are tested using automated testing facilities provided by Golang. 
  * When testing the application uses stubbing of the third-party endpoints to ensure test reliability (removing dependency on external services).
  * Testing includes testing of handlers using the httptest package, as well as unit tests. 
  * Test coverage of TODO ----------------- !!!!!!!!!!! percent.
* Repeated invocations for a given country and date are cached on firebase to minimise invocation on the third-party libraries. These are deleted if the cached requests are older than a constant that can be set by the user. The default value is 4 hours. 

Allocation of tasks: 
| Functionality | Name | 
| ----- | ------ |
|All endpoints | Kevin |
|Error handling | Mostly Raphael, also Sondre and Kevin  |
|Webhooks | Kevin |
|Restcountries interaction | Mostly Raphael, also Sondre |
|Cache | Kevin |
|Firestore setup and connectivity | Kevin |
|Http testing | Torje |
|Unit testing | Sondre, Torje |
|Stubbing of third-party services| Sondre |
|Openstack deployment | Raphael  |
|Dockerfile and docker compose | Raphael  |
|Readme.md | Kevin, Raphael |
|Debugging | Everyone |

# Running the assignment


# Running the Assignment

## Openstack Instance

The easiest way to access the assignment is by using the official deployment.

1. Ensure you are connected to the NTNU network, either physically or through a VPN tunnel.
2. Use the following URL to access the API: http://10.212.171.254:80

## Docker

The application can also be run locally as a Docker container if Docker is installed on your computer. To run the project locally, you must first set up a database.

1. Sign up for https://firebase.google.com/ and create a new Firestore database.
2. Pull the service account JSON from the project settings tab on the Firestore web UI, and place it within the ".secrets" folder of the repository. If this folder doesn't exist, create one in the root of the repository.
3. Rename the file to "production_credentials.json".
4. Once the previous steps are complete, run the command `docker-compose up -d` while having the root of the repository as your working directory (on Linux).
5. After a few minutes of loading, you will see that six containers have been created, whereas one of them is immediately deleted. At this point, the installation is complete. A new collection containing some country information will have been added to your Firestore database. Navigate to "localhost:80" to access the API.

## Manual Compile

The application can also be run locally using the "go run" command. This method is not recommended, as it requires you to navigate into the source code for it to be able to run.

1. As with the Docker method, you must first create a new Firestore database before continuing.
2. Create a new folder within the repository root named "credentials". Place the Firebase service account JSON within this folder.
3. Rename the service account file to "production_credentials.json".
4. Open the file `/utils/constants/constants.go`. Change the `CREDENTIALS_FILE` constant variable to equal `./credentials/production_credentials.json`.
5. In the command line interface, enter the following commands: `go get assignment2/utils/db`, `go get assignment2/handlers`, and `go get github.com/stretchr/testify/assert`. This will install external libraries necessary to run the application.
6. Enter `go run ./cmd/main.go` in the command line interface to compile and run the application.
7. The API is now available at "localhost:8080", or a different port if specified in the environment.


# How to use the API

The usage of this service should follow the following specifications for schemas (or syntax) of requests.

The responses will follow the specifications bellow for reponse body, as well as method and status code. 

## Endpoints

The web service has four resource root paths: 

```
/energy/v1/renewables/current
/energy/v1/renewables/history
/energy/v1/notifications/
/energy/v1/status/
```

The specification has the following conventions for placeholders:

* {value} - *mandatory* value
* {value?} - *optional* value
* {?key=value} - *mandatory* parameter (key-value pair)
* {?key=value?} - *optional* parameter (key-value pair)

## Current percentage of renewables

This endpoint returns the latest percentages of renewables in the energy mix.

### - Request

```
Method: GET
Path: /energy/v1/renewables/current/{country?}{?neighbours=bool?}{?sortByValue=bool?}
```

`{country?}` refers to an optional country identifier, either a 3-letter code **or** the name of the country.

`{?neighbours=bool?}` refers to an optional parameter indicating whether neighbouring countries' values should be shown. Will be ignored if no country is given. 

`{?sortByValue=bool?}` refers to an optional parameter indicating whether the output will be sort by percentage value (e.g., `?sortByValue=true`).

Example request:
* ```/energy/v1/renewables/current/nor```
* ```/energy/v1/renewables/current/norway?neighbours=true```
* ```/energy/v1/renewables/current/sweden?neighbours=true&sortByValue=true```
* ```/energy/v1/renewables/current/```
* ```/energy/v1/renewables/current/?sortByValue=true```
### - Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise indicating wether the request is illegal or there has been a server error.

Body (Exemplary message based on schema) - *with* country code:
```
{
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2021",
    "percentage": 71.558365
}
```

Body (Exemplary message based on schema) - *with* country code *and* neighbour parameter activated:
```
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "2021",
        "percentage": 71.558365
    },
    {
        "name": "Finland",
        "isoCode": "FIN",
        "year": "2021",
        "percentage": 34.61129
    },
    {
        "name": "Russia",
        "isoCode": "RUS",
        "year": "2021",
        "percentage": 6.6202893
    },
    {
        "name": "Sweden",
        "isoCode": "SWE",
        "year": "2021",
        "percentage": 50.924007
    }
]
```

Body (Exemplary message based on schema) - *without* country code (returns all countries)

```
[
    {
        "name": "Algeria",
        "isoCode": "DZA",
        "year": "2021",
        "percentage": 0.26136735
    },
    {
        "name": "Argentina",
        "isoCode": "ARG",
        "year": "2021",
        "percentage": 11.329249
    },
    {
        "name": "Australia",
        "isoCode": "AUS",
        "year": "2021",
        "percentage": 12.933532
    },
    ...
]
```

Body (Exemplary message based on schema) - *without* country code *and* sortByValue parameter activated: (returns all countries sorted by percentage value)

```
[
    {
        "name": "Iceland",
        "isoCode": "ISL",
        "year": "2021",
        "percentage": 86.874535
    },
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "2021",
        "percentage": 71.558365
    },
    {
        "name": "Sweden",
        "isoCode": "SWE",
        "year": "2021",
        "percentage": 50.924007
    },
    ...
]
```

## Historical percentages of renewables

This endpoint focuses on returning historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

### - Request

```
Method: GET
Path: /energy/v1/renewables/history/{country?}{?begin=year}{?end=year?}{?neighbours=bool?}{?sortByValue=bool?}{?mean=bool?}
```

`{country?}` refers to an optional country identifier, either a 3-letter code **or** the name of the country.

`{?begin=year}` refers to an optional parameter indicating the earliest year of data the output will contain. No earlier years, and all laters years will be included (except if defined otherwise by the end parameter). If the output is mean percentage, the mean value will only be calculated from data later than this value.   

`{?end=year}` refers to an optional parameter indicating the lastest year of data the output will contain. No later years, and all previous years will be included (except if defined otherwise by the begin parameter). If the output is mean percentage, the mean value will only be calculated from data earlier than this value.  

`{?neighbours=bool?}` refers to an optional parameter indicating whether neighbouring countries' values should be shown. Will be ignored if no country is given. 

`{?sortByValue=bool?}` refers to an optional parameter indicating whether the output will be sort by percentage value (e.g., `?sortByValue=true`).

 `{?mean=bool?}` refers to an optional parameter indicating whether the output will be the mean value instead of data for each year. Will be ignored if no country is given, as this will always return mean value. 


Example request: 
* ```/energy/v1/renewables/history/nor```
* ```/energy/v1/renewables/history/norway?begin=2000```
* ```/energy/v1/renewables/history/NOR?begin=2010&end=2020&neighbours=true&sortByValue=true```
* ```/energy/v1/renewables/history/NOR?begin=1990&mean=true```
* ```/energy/v1/renewables/history/```
* ```/energy/v1/renewables/history/end=1975```
* ```/energy/v1/renewables/history/?sortByValue=true```

### - Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise indicating wether the request is illegal or there has been a server error.

Body (Exemplary message based on schema) - *with* country code:
```
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1965",
        "percentage": 67.87996
    },
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1966",
        "percentage": 65.3991
    },
    ...
]
```

Body (Exemplary message based on schema) - *with* country code, and mean, neighbours, and sortByValue set to true:
```
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "percentage": 68.01918892982457
    },
    {
        "name": "Sweden",
        "isoCode": "SWE",
        "percentage": 33.97086068421053
    },
    {
        "name": "Finland",
        "isoCode": "FIN",
        "percentage": 18.825984771929832
    },
    ...
]
```

Body (Exemplary message based on schema) - *without* country code (returns mean percentages for all countries):
```
[
    {
        "name": "United Arab Emirates",
        "isoCode": "ARE",
        "percentage": 0.0444305504
    },
    {
        "name": "Argentina",
        "isoCode": "ARG",
        "percentage": 9.131337212280702
    },
    {
        "name": "Australia",
        "isoCode": "AUS",
        "percentage": 5.3000481596491245
    },
    ...
]
```

## Notification Endpoint

Users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries (or any country) is invoked, where the minimum frequency can be specified. If specified, a webhook can only be triggered at the specified year. Users can register multiple webhooks. The registrations will be stored until explicitly deleted. 

### Registration of Webhook

### - Request

```
Method: POST
Path: /energy/v1/notifications/
```

* Content type: `application/json`

The body contains 
 * the URL to be triggered upon event (the service that should be invoked)
 * the country for which the trigger applies (if empty, it applies to any invocation)
 * the number of invocations after which a notification is triggered (it should re-occur every *number of invocations*, i.e., if 5 is specified, it should occur after 5, 10, 15 invocation, and so on, unless the webhook is deleted).
 * an optional value "year" which specify for which year the trigger applies (if empty it applies to any year)

Body (Exemplary message based on schema):
```
{
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}
```

Body (Exemplary message based on schema) with year:
```
{
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5,
   "year": 2000
}
```
### - Response

The response contains the ID for the registration that can be used to see detail information or to delete the webhook registration. The format of the ID is a unique randomly generated 16 character string.

* Content type: `application/json`
* Status code: 201 Status created if everything is OK, appropriate error code otherwise indicating wether the request is illegal or there has been a server error.

Body (Exemplary message based on schema):
```
{
    "webhook_id": "BOlOomFOeiKvZhVD"
}
```

### Deletion of Webhook

### - Request

```
Method: DELETE
Path: /energy/v1/notifications/{id}
```

* {id} is the ID returned during the webhook registration

### - Response

* Status code: 204 No content if everything is OK, appropriate error code otherwise indicating wether the request is illegal or there has been a server error.


### View registered webhook

### - Request

```
Method: GET
Path: /energy/v1/notifications/{id}
```
* `{id}` is the ID for the webhook registration

### - Response

The response is similar to the POST request body, but further includes the ID assigned by the server upon adding the webhook.

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise indicating wether the request is illegal or there has been a server error.

Body (Exemplary message based on schema):
```
{
   "webhook_id": "BOlOomFOeiKvZhVD",
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}

```
Body (Exemplary message based on schema) with no country specified and year:
```
{
   "webhook_id": "QDzPVIWGuZkfueZx",
   "url": "https://localhost:8081/client/",
   "country": "ANY",
   "calls": 2,
   "year": 2020
}

```

### View all registered webhooks

### - Request

```
Method: GET
Path: /energy/v1/notifications/
```

### - Response

The response is a collection of all registered webhooks.

* Content type: `application/json`

Body (Exemplary message based on schema):
```
[
   {
      "webhook_id": "BOlOomFOeiKvZhVD",
      "url": "https://localhost:8080/client/",
      "country": "NOR",
      "calls": 5
   },
   {
      "webhook_id": "QDzPVIWGuZkfueZx",
      "url": "https://localhost:8081/client/",
      "country": "ANY",
      "calls": 2,
      "year": 2020
    },
   ...
]
```

### Webhook Invocation (upon trigger)

When a webhook is triggered, it sends information as follows. Where multiple webhooks are triggered, the information is sent separately. 

```
Method: POST
Path: <url specified in the corresponding webhook registration>
```

* Content type: `application/json`

Body (Exemplary message based on schema):
```
{
   "webhook_id": "BOlOomFOeiKvZhVD",
   "country": "Norway",
   "calls": 10
}
```

Body (Exemplary message based on schema) where no country is specified:
```
{
   "webhook_id": "QfwLosaJKVANmUJk",
   "calls": 4
}
```

Body (Exemplary message based on schema) when year is specified:
```
{
   "webhook_id": "ScFdJSpMVIMsXznf",
   "country": "Sweden",
   "calls": 8,
   "year": 2020
}
```
* Note: `calls` show the number of invocations, not the number specified as part of the webhook registration (i.e. the actual invocation upon which the webhook is triggered).

## Status Endpoint

The status interface indicates the availability of all individual services this service depends on. The reporting occurs based on status codes returned by the dependent services. The status interface further provides information about the number of registered webhooks and the uptime of the service.

### - Request

```
Method: GET
Path: energy/v1/status/
```

### - Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. 

Body:
```
{
   "countries_api": "<http status code for *REST Countries API*>",
   "notification_db": "<http status code for *Notification DB* in Firebase>",
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```

