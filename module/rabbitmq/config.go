package rabbitmq

import (
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/metricbeat/mb"
)

type Config struct {
	Hosts    []string           `config:"hosts" validate:"nonzero,required"`
	Username string             `config:"username"`
	Password string             `config:"password"`
	TLS      *outputs.TLSConfig `config:"ssl"`
}

var defaultConfig = Config{
	Hosts:    []string{"127.0.0.1"},
	Username: "guest",
	Password: "guest",
	TLS:      nil,
}

func CreateClient(base mb.BaseMetricSet) (Api, error) {

	config := defaultConfig
	err := base.Module().UnpackConfig(&config)

	if err != nil {
		return nil, err
	}

	return NewClient(base.Host(), config.Username, config.Password, config.TLS)
}
