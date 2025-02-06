package config

type DialogServiceConfig struct {
	DefaultPageSize int  `yaml:"default_page_size"`
	MaxPageSize     int  `yaml:"max_page_size"`
	UseInmem        bool `yaml:"use_inmem"`
}
