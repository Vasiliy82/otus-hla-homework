package config

type RouteConfig struct {
	Path     string           `yaml:"path"`
	Methods  []string         `yaml:"methods"`
	Services []*ServiceConfig `yaml:"services"`
}

type ServiceConfig struct {
	ServiceName       string   `yaml:"service_name"`
	URL               string   `yaml:"url"`
	SupportedVersions []string `yaml:"supported_versions"`
}
