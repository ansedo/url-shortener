package config

const (
	SiteScheme = "http://"
	SiteHost   = "localhost"
	SitePort   = ":8080"

	SiteAddress = SiteScheme + SiteHost + SitePort

	RequestNotAllowedError = "this request is not allowed"
)
