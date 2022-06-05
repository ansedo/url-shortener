package config

type Config struct {
	SiteScheme  string
	SiteHost    string
	SitePort    string
	SiteAddress string

	RequestNotAllowedError string
}

func NewConfig() *Config {
	cfg := &Config{
		SiteScheme: "https://",
		SiteHost:   "localhost",
		SitePort:   ":8080",

		RequestNotAllowedError: "this request is not allowed",
	}
	cfg.SiteAddress = cfg.SiteScheme + cfg.SiteHost + cfg.SitePort
	return cfg
}
