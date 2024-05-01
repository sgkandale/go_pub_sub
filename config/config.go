package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GrcpConfig      ServerConfig
	WebsocketConfig ServerConfig
}

type ServerConfig struct {
	Use  bool
	Port int
}

type FileConfig struct {
	// grpc config
	UseGrpc  bool `yaml:"use_grpc" env:"use_grpc"`
	GrpcPort int  `yaml:"grpc_port" env:"grpc_port"`

	// websocket config
	UseWebsocket  bool `yaml:"use_websocket" env:"use_websocket"`
	WebsocketPort int  `yaml:"websocket_port" env:"websocket_port"`
}

func Read(fileName string) *Config {
	var cfg FileConfig
	log.Printf("[INFO] reading config from file : %s", fileName)
	err := cleanenv.ReadConfig(fileName, &cfg)
	if err != nil {
		log.Printf("[ERROR] reading config file %s : %s", fileName, err)
		log.Print("[INFO] reading config from env variables")
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatal("[ERROR] reading config using env variables : ", err)
		}
	}
	log.Print("[INFO] config reading done")

	respCfg := &Config{}

	// if grpc is enabled
	if cfg.UseGrpc {
		log.Print("[INFO] grpc server is enabled in config")
		respCfg.GrcpConfig = ServerConfig{
			Use:  true,
			Port: cfg.GrpcPort,
		}
	} else {
		log.Print("[INFO] grpc server is not enabled in config")
	}

	// if websocket is enabled
	if cfg.UseWebsocket {
		log.Print("[INFO] websocket server is enabled in config")
		respCfg.WebsocketConfig = ServerConfig{
			Use:  true,
			Port: cfg.WebsocketPort,
		}
	} else {
		log.Print("[INFO] websocket server is not enabled in config")
	}

	return respCfg
}
