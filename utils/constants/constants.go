package constants

const VERSION = "v1" // Service version

// Endpoint paths

const DEFAULT_PATH = "/"                                      // Default path
const SERVICE_PATH = "/energy/" + VERSION                     // Service path
const RENEWABLES_PATH = SERVICE_PATH + "/renewables"          // Renewables path
const RENEWABLES_CURRENT_PATH = RENEWABLES_PATH + "/current/" // Renewables current path
const RENEWABLES_HISTORY_PATH = RENEWABLES_PATH + "/history/" // Renewables history path
const NOTIFICATION_PATH = SERVICE_PATH + "/notification/"     // Notification path
const STATUS_PATH = SERVICE_PATH + "/status"                  // Status path

// Content type

const CONT_TYPE_JSON = "application/json" // Content type JSON

// Country API

const COUNTRIES_API_URL = "http://129.241.150.113:8080" // URL to countries API
const COUNTRY_NAME_SEARCH_PATH = "/v3.1/name/"          // Path to search for country name
const COUNTRY_CODE_SEARCH_PATH = "/v3.1/alpha/"         // Path to search for country code
const USED_COUNTRY_CODE = "cca3"                        // Country code used in response from countries API

// Years for renewables database

const OLDEST_YEAR_DB = 1965 // Oldest year in database
const LATEST_YEAR_DB = 2021 // Latest year in database

// Firestore constants

const RENEWABLES_COLLECTION = "renewables" // Name of renewables collection
const WEBHOOKS_COLLECTION = "webhooks"     // Name of webhooks collection
const CACHE_COLLECTION = "cache"           // Name of cache collection

// Name of files
const RENEWABLES_CSV_FILE = "./res/renewable-share-energy.csv"             // Path to CSV file
const CREDENTIALS_FILE = "/credentials/production_credentials.json"        // Path to credentials file
const CREDENTIALS_FILE_TESTING = "../credentials/testing_credentials.json" // Path to credentials file for testing
const RESTCOUNTRIES_MOCK = "./res/restcountries-mock.json"                 // Path to mock file for restcountries API

// Webhooks

const WEBHOOK_ID_LENGTH = 16 // Length of webhook ID

// Cache

const MAX_CACHE_AGE_IN_HOURS = 4 // Max age of cache in hours

// Default error responses

const DEFAULT500 = "There has been an internal server error. Please try again later."                                // Default 500 error message
const DEFAULT503 = "Service is currently unavailable. Please try again later."                                       // Default 503 error message
const DEFAULT504 = "Service is currently down due to the failure of an external dependency. Please try again later." // Default 504 error message
