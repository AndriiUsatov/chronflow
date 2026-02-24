package config

import (
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type JanitorConfig struct {
	PGURI        string `mapstructure:"pg_uri" validate:"required"`
	PGPort       int    `mapstructure:"pg_port" validate:"required"`
	PGUser       string `mapstructure:"pg_user" validate:"required"`
	PGPwd        string `mapstructure:"pg_pwd" validate:"required"`
	PGTaskDB     string `mapstructure:"pg_task_db" validate:"required"`
	PGSSLMode    string `mapstructure:"pg_ssl_mode" validate:"required"`
	PGTaskSchema string `mapstructure:"pg_task_schema" validate:"required"`
	PGTaskTable  string `mapstructure:"pg_task_table" validate:"required"`

	JanitorHeartbeatPort string `mapstructure:"janitor_heartbeat_port" validate:"required"`
}

func LoadApiConfig(path string) (cfg JanitorConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("CHRONFLOW")
	viper.AutomaticEnv()

	viper.SetDefault("pg_uri", "")
	viper.SetDefault("pg_port", 0)
	viper.SetDefault("pg_user", "")
	viper.SetDefault("pg_pwd", "")
	viper.SetDefault("pg_task_db", "")
	viper.SetDefault("pg_ssl_mode", "")
	viper.SetDefault("pg_task_schema", "")
	viper.SetDefault("pg_task_table", "")
	viper.SetDefault("janitor_heartbeat_port", "")

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
