package client

var HOST = "https://eu.onetimesecret.com"
var BASE_URI = "https://eu.onetimesecret.com/api"
const API_VERSION = "v1"

var ENDPOINTS = map[string]string{
	"status":      "status",
	"share":       "share",
	"generate":    "generate",
	"getsecret":   "secret",
	"getmetadata": "private",
	"burn":        "private",
	"getrecent":   "private",
}
