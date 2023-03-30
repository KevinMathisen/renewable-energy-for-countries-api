package constants

// Service Version
const VERSION = "v1"

// Endpoint Paths
const DEFAULT_PATH = "/"
const SERVICE_PATH = "/energy/" + VERSION
const RENEWABLES_PATH = SERVICE_PATH + "/renewables"
const RENEWABLES_CURRENT_PATH = RENEWABLES_PATH + "/current"
const RENEWABLES_HISTORY_PATH = RENEWABLES_PATH + "/history"
const NOTIFICATION_PATH = SERVICE_PATH + "/notification"
const STATUS_PATH = SERVICE_PATH + "/status"

// Content type
const CONT_TYPE_JSON = "application/json"
