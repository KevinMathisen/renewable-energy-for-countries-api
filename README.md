# Assignment2

# Overview

This service is a REST web application in Golang that provides the client with the ability to retrieve information about developments related to renewable energy production for and across countries. It uses an existing webservice, [restcountries](http://129.241.150.113:8080) for translation isoCodes and names, as well as for looking up border countries. It also uses a firestore database for saving renewables data, as well as webhooks registered, and a cache.

The service allows for notification registration using webhooks, which are invoked based on requests to specific countries (and years if specified).  

The application is dockerized and deployed using an IaaS system called Openstack on NTNUs instance called SkyHigh. See Running the assignment/Openstack instance 


The REST web service:
* *REST Countries API*. Endpoint: http://129.241.150.113:8080/v3.1 (Documentation: http://129.241.150.113:8080/)

Dataset used for Renewables:
* [*Renewable Energy Dataset*](https://drive.google.com/file/d/18G470pU2NRniDfAYJ27XgHyrWOThP__p/view?usp=sharing) (Authors: Hannah Ritchie, Max Roser and Pablo Rosado (2022) - "Energy". Published online at OurWorldInData.org. Retrieved from: https://ourworldindata.org/energy

The dataset reports on percentage of renewable energy in the country's energy mix over time. 

# Running the assignment

## Openstack instance

The Openstack deployment is located at http://10.212.171.254:8080

## Locally
To run the project locally, you must first set up a database. Sign up to https://firebase.google.com/, and create a new firestone database. From the project settings tab, pull the service account credentials.


If you downloaded the repository you can run it in one of two ways, as long as you have golang installed on your computer.
- In the repository run the command `go run ./cmds/main.go`
- In the repository, first run the command `go build ./cmds`, which will create a file calles *cmds.exe*. Running this file will start the application.

When running the application locally it will be available on `localhost:8080`, or a different port if specified in environment. 

## Docker

The application may also be run locally as a docker container if docker is installed on the computer. To run the project locally, you must first set up a database. Sign up to https://firebase.google.com/, and create a new firestone database. Then, pull the service account json from the project settings tab on the firstone web-ui, and place it within the ".secrets" folder of the repository. If this folder doesn't exist, create one in the root of the repository. When this is complete, 


# How to use the API

The usage of this service should follow the following specifications for schemas (or syntax) of requests.

The responses will follow the specifications bellow for reponse body, as well as method and status code. 

## Endpoints

Your web service has four resource root paths: 

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

The initial endpoint returns the latest percentages of renewables in the energy mix.

### - Request

```
Method: GET
Path: /energy/v1/renewables/current/{country?}
```

`{country?}` refers to an optional country identifier, either a 3-letter code or the name of the country.

`{?neighbours=bool?}` refers to an optional parameter indicating whether neighbouring countries' values should be shown. 

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

## Historical percentages of renewables

The initial endpoint focuses on returning historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

### - Request

```
Method: GET
Path: /energy/v1/renewables/history/{country?}{?begin=year&end=year?}
```

```{country?}``` refers to an optional country 3-letter code.

Example request: ```/energy/v1/renewables/history/nor```

### - Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. Ensure to deal with errors gracefully.

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
* Where `{?begin=year&end=year?}` is specified (e.g., `?begin=1960&end=1970`), only these years should be shown.

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

* **Advanced Tasks:** 
  * Consider selective use of only `begin` or `end` as single parameter (e.g., `?begin=1980` should only consider data from 1980 onwards; `?end=1980` should consider values from the first time entry until 1980 only).
  * Extend the history for all countries with a time constraint. That means, where `{?begin=year&end=year?}` is specified (e.g., `?begin=1960&end=1970`), only mean values for these years should be calculated (not for all years).
  * Implement additional optional parameter `{?sortByValue=bool?}` to support sorting of output by percentage value (e.g., `?sortByValue=true`).
  * be creative - can you think of other/alternative useful options? Whatever you implement in addition, remember to document this.

## Notification Endpoint

As an additional feature, users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries (or any country) is invoked, where the minimum frequency can be specified. Users can register multiple webhooks. The registrations should survive a service restart (i.e., be persistent using a Firebase DB as backend).

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

Body (Exemplary message based on schema):
```
{
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
}
```
### - Response

The response contains the ID for the registration that can be used to see detail information or to delete the webhook registration. The format of the ID is not prescribed, as long it is unique. Consider best practices for determining IDs.

* Content type: `application/json`
* Status code: Choose an appropriate status code

Body (Exemplary message based on schema):
```
{
    "webhook_id": "OIdksUDwveiwe"
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

Implement the response according to best practices.

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

Body (Exemplary message based on schema):
```
{
   "webhook_id": "OIdksUDwveiwe",
   "url": "https://localhost:8080/client/",
   "country": "NOR",
   "calls": 5
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
      "webhook_id": "OIdksUDwveiwe",
      "url": "https://localhost:8080/client/",
      "country": "NOR",
      "calls": 5
   },
   {
      "webhook_id": "DiSoisivucios",
      "url": "https://localhost:8081/anotherClient/",
      "country": "SWE",
      "calls": 2
   },
   ...
]
```

### Webhook Invocation (upon trigger)

When a webhook is triggered, it should send information as follows. Where multiple webhooks are triggered, the information should be sent separately (i.e., one notification per triggered webhook). Note that for testing purposes, this will require you to set up another service that is able to receive the invocation. During the development, consider using https://webhook.site/ initially.

```
Method: POST
Path: <url specified in the corresponding webhook registration>
```

* Content type: `application/json`

Body (Exemplary message based on schema):
```
{
   "webhook_id": "OIdksUDwveiwe",
   "country": "Norway",
   "calls": 10
}
```
* Note: `calls` should show the number of invocations, not the number specified as part of the webhook registration (i.e., not 5, but the actual invocation upon which the webhook is triggered).

* **Advanced Task:** Consider supporting other event types you can think of.

## Status Endpoint

The status interface indicates the availability of all individual services this service depends on. These can be more services than the ones specified above (if you considered the advanced tasks). If you include more, you can specify additional keys with the suffix `api`. The reporting occurs based on status codes returned by the dependent services. The status interface further provides information about the number of registered webhooks (more details is provided in the next section), and the uptime of the service.

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
   ...
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```

Note: `<some value>` indicates placeholders for values to be populated by the service as described for the corresponding values. Feel free to extend the output with information you deem useful to assess the status of your service.

# Additional requirements

* All endpoints should be *tested using automated testing facilities provided by Golang*. 
  * This includes the stubbing of the third-party endpoints to ensure test reliability (removing dependency on external services).
  * Include the testing of handlers using the httptest package. Your code should be structured to support this. 
  * Try to maximize test coverage as reported by Golang.
* Repeated invocations for a given country and date should be cached to minimise invocation on the third-party libraries. Use Firebase for this purpose.
  * **Advanced Task**: Implement purging of cached information for requests older than a given number of hours/days.

# Deployment

The service is to be deployed on an IaaS solution OpenStack using Docker (to be discussed in class). You will need to provide the URL to the deployed service as part of the submission, in addition the source repository.

# Notes

* Feel free to introduce additional endpoints to support the development and debugging.
* Where specification details are missing (but you can infer those), operate based on best practices and document it accordingly.
* Where information is unclear, get in touch with teaching staff for clarification. Where needed, the assignment information will be updated accordingly (and the updated information will be highlighted as **UPDATE**).

# General Aspects

## Professionalism

As indicated during the initial sessions, ensure you work with professionalism in mind (see Course Rules). In addition to professionalism, you are at liberty to introduce further features into your service, as long it does not break the specification given above. 

## Workspace environment

Please work in the provided workspace environment (see [here](Rules-&-Conventions/Workspace-Conventions) - lodge an issue if you have trouble accessing it) for your user and create a project `assignment-2` in this workspace. All group members share *one* repository; it does not matter in whose workspace folder it lies.

## Rate limits reminder

As mentioned above, be sensitive to rate limits of external services. This has proven very important, given the large number of projects (and hence invocations) on the third-party services.

## Resources

The course repository provides a range of example projects for various features discussed throughout the lecture sessions. Feel free to borrow from those projects, or use them to understand a concept you are struggling with (e.g., learning the use of Firestore). 

## Third-party libraries

Be deliberative about using third-party libraries (Don't just do it because someone did it on StackOverflow). While those libraries often allow for convenience and functionality you would otherwise need to reimplement, they can also mean the "import" of technological debt, especially if you were to think about maintainability. So, be very clear *why* you want to use the library (lack of functionality in standard packages, convenience, etc.). Note that it may be challenging for us to provide the necessary support, especially if the library is rather specialised (we will rely on the same resources available to you). As the assignment is designed, you will only need Golang standard API functionality, alongside Firebase/store functionality as a third-party dependency.

# Submission

The **assignment is a group assignment**. Ensure that the group allocations specified in the submission system are correct at the deadline. The **submission deadline** is provided on the [course wiki page](Home#deadlines). Extensions to the deadline are handled according to the [Course Rules](Rules-&-Conventions/Course%20Rules). 

As part of the submission you will need to provide:
* a link to your code repository (ensure it is `internal` at that stage)
* a link to the deployed service

In addition, we will provide you with an option to clarify aspects of your submission (e.g., checklist of features, elaboration on aspects that don't quite work, or additional features).

The submission occurs via our [Submission System](Guides/Submission System). Early submission is explicitly encouraged - you can change it (or even withdraw) any time before the deadline.

# Peer Review

After the submission deadline, there will be a separate deadline during which you will review other groups' submissions. To do this the system provides you with a checklist of aspects to assess. You will need to review *at least two submissions* to meet the mandatory requirements of peer review, but you can review as many submissions as you like, which counts towards your participation mark for the course. The peer-review deadline is indicated on the [course wiki page](Home#deadlines).
