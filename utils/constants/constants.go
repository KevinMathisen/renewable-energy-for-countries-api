package constants

// Service version
const VERSION = "v1"

// Endpoint paths
const DEFAULT_PATH = "/"
const SERVICE_PATH = "/energy/" + VERSION
const RENEWABLES_PATH = SERVICE_PATH + "/renewables"
const RENEWABLES_CURRENT_PATH = RENEWABLES_PATH + "/current/"
const RENEWABLES_HISTORY_PATH = RENEWABLES_PATH + "/history/"
const NOTIFICATION_PATH = SERVICE_PATH + "/notification"
const STATUS_PATH = SERVICE_PATH + "/status"

// Content type
const CONT_TYPE_JSON = "application/json"

// Country API
const COUNTRIES_API_URL = "http://129.241.150.113:8080"
const COUNTRY_NAME_SEARCH_PATH = "/v3.1/name/"
const ISO_SEARCH_PATH = "/v3.1/alpha/"

// Years for renewables database
const OLDEST_YEAR_DB = 1965
const LATEST_YEAR_DB = 2021

// Firestore constants
const RENEWABLES_COLLECTION = "renewables"
const WEBHOOKS_COLLECTION = "webhooks"
const CACHE_COLLECTION = "cache"

// Name of files
const RENEWABLES_CSV_FILE = "./renewable-share-energy.csv"
const CREDENTIALS_FILE = "./assignment2-prog2005-service-account.json"

// Webhooks
const WEBHOOK_ID_LENGTH = 16
const WEBHOOK_COLLECTIONNAME = "webhooks"
