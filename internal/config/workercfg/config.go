package config

import (
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	NatsURL                               string `mapstructure:"nats_url" validate:"required"`
	NatsTaskStream                        string `mapstructure:"nats_task_stream" validate:"required"`
	NatsTaskToProcessSubject              string `mapstructure:"nats_task_to_process_subject" validate:"required"`
	NatsTaskToProcessDurable              string `mapstructure:"nats_task_to_process_durable" validate:"required"`
	GrpcTaskUpdateServerTransportProtocol string `mapstructure:"grpc_task_update_server_transp" validate:"required"`
	GrpcTaskUpdateServerURL               string `mapstructure:"grpc_task_update_server_url" validate:"required"`
	GrpcTaskUpdateServerPort              int    `mapstructure:"grpc_task_update_server_port" validate:"required"`
}

func LoadWorkerConfig(path string) (cfg WorkerConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("CHRONFLOW")
	viper.AutomaticEnv()

	viper.SetDefault("nats_url", "")
	viper.SetDefault("nats_task_stream", "")
	viper.SetDefault("nats_task_to_process_subject", "")
	viper.SetDefault("nats_task_to_process_durable", "")
	viper.SetDefault("grpc_task_update_server_transp", "")
	viper.SetDefault("grpc_task_update_server_url", "")
	viper.SetDefault("grpc_task_update_server_port", 0)

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}

	if err = validator.New().Struct(&cfg); err != nil {
		return
	}

	return
}
