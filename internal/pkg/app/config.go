package app

import (
	"github.com/nil-go/konf"
	"github.com/nil-go/konf/provider/env"
	"github.com/nil-go/konf/provider/file"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging struct {
		Level int
	}
	Web  WebConfig
	GRPC GrpcConfig
	DB   struct {
		DriverName       string
		ConnectionString string
	}
}

func ReadLocalConfig(configPath string) (Config, error) {
	var config konf.Config

	err := config.Load(file.New(configPath, file.WithUnmarshal(yaml.Unmarshal)))
	if err != nil {
		return Config{}, err
	}

	err = config.Load(env.New())
	if err != nil {
		return Config{}, err
	}

	var res Config

	err = config.Unmarshal("", &res)
	if err != nil {
		return Config{}, err
	}

	return res, nil
}
