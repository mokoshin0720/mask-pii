package config

import (
	"github.com/caarlos0/env/v6"
	"go.uber.org/multierr"
)

var Config config

type config struct {
	GcpApiKey string `env:"GCP_API_KEY,notEmpty"`
	InputText string `env:"INPUT_TEXT,notEmpty"`
}

func Init() error {
	var err error
	if enverr := env.Parse(&Config); enverr != nil {
		err = multierr.Append(err, enverr)
	}

	return err
}
