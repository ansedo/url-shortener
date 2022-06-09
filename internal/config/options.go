package config

type Option func(c *config)

func WithSiteScheme(siteScheme string) Option {
	return func(cfg *config) {
		cfg.SiteScheme = siteScheme
	}
}

func WithSiteHost(siteHost string) Option {
	return func(cfg *config) {
		cfg.SiteHost = siteHost
	}
}

func WithSitePort(sitePort string) Option {
	return func(cfg *config) {
		cfg.SitePort = sitePort
	}
}

func WithRequestNotAllowedError(requestNotAllowedError string) Option {
	return func(cfg *config) {
		cfg.RequestNotAllowedError = requestNotAllowedError
	}
}
