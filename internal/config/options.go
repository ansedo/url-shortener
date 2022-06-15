package config

type Option func(c *config)

func WithBaseURL(baseURL string) Option {
	return func(cfg *config) {
		cfg.BaseURL = baseURL
	}
}

func WithServerAddress(serverAddress string) Option {
	return func(cfg *config) {
		cfg.ServerAddress = serverAddress
	}
}

func WithFileStoragePath(fileStoragePath string) Option {
	return func(cfg *config) {
		cfg.FileStoragePath = fileStoragePath
	}
}
