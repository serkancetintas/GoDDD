package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	configPath = "../common/resources"
	configType = "yaml"
)

func NewConfigurationManager() IConfigurationManager {
	env := os.Getenv("PROFILE")
	if env == "" {
		env = "local"
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigType(configType)
	applicationConfig := readApplicationConf(env)
	queuesConfig := readQueuesConfig()
	return &ConfigurationManager{applicationConfig: applicationConfig, queuesConfig: queuesConfig}
}

type IConfigurationManager interface {
	GetRabbitConfig() RabbitConfig
	GetQueuesConfig() QueuesConfig
}

type ConfigurationManager struct {
	applicationConfig ApplicationConfig
	queuesConfig      QueuesConfig
}

func (configurationManager *ConfigurationManager) GetRabbitConfig() RabbitConfig {
	return configurationManager.applicationConfig.Rabbit
}

func (configurationManager *ConfigurationManager) GetQueuesConfig() QueuesConfig {
	return configurationManager.queuesConfig
}

func readQueuesConfig() QueuesConfig {
	viper.SetConfigName("rabbit-queue")
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		log.Panicf("Couldn't load queues configuration, cannot start. Error details: %s", readConfigErr.Error())
	}
	var conf QueuesConfig
	c := viper.Sub("queue")
	unMarshalErr := c.Unmarshal(&conf)
	unMarshalSubErr := c.Unmarshal(&conf)
	if unMarshalErr != nil {
		log.Panicf("Configuration cannot deserialize. Terminating. Error details: %s", unMarshalErr.Error())
	}
	if unMarshalSubErr != nil {
		log.Panicf("Configuration cannot deserialize. Terminating. Error details: %s", unMarshalSubErr.Error())
	}
	logrus.WithField("configuration", conf).Debug("Configuration changed")
	return conf
}

func readApplicationConf(env string) ApplicationConfig {
	viper.SetConfigName("application")
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		log.Panicf("Couldn't load application configuration, cannot start. Error details: %s", readConfigErr.Error())
	}

	var conf ApplicationConfig
	c := viper.Sub(env)
	unMarshalErr := c.Unmarshal(&conf)
	if unMarshalErr != nil {
		log.Panicf("Configuration cannot deserialize. Terminating. Error details: %s", unMarshalErr.Error())
	}

	unMarshalSubErr := c.Unmarshal(&conf)
	if unMarshalSubErr != nil {
		log.Panicf("Configuration cannot deserialize. Terminating. Error details: %s", unMarshalSubErr.Error())
	}

	logrus.WithField("configuration", conf).Debug("Configuration changed")
	return conf
}
