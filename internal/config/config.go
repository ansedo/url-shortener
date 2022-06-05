package config

type Config struct {
	SiteScheme  string
	SiteHost    string
	SitePort    string
	SiteAddress string

	RequestNotAllowedError string
}

func NewConfig() *Config {
	config := &Config{
		SiteScheme: "https://",
		SiteHost:   "localhost",
		SitePort:   ":8080",

		RequestNotAllowedError: "this request is not allowed",
	}
	config.SiteAddress = config.SiteScheme + config.SiteHost + config.SitePort
	return config
}
