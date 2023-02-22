package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	IsDebug bool
	Server  Server
	Mysql   Mysql
	MinIO   MinIO
	Redis   Redis
	MongoDB MongoDB
	Nats    Nats
}

type Server struct {
	Address  string
	HttpPort int
}

type Mysql struct {
	Uri         string
	MaxConn     int
	MaxIdel     int
	MaxLifeTime int
}

type MinIO struct {
	DoMain       string
	EndPoint     string
	AccessSecret string
	AccessKey    string
	BucketName   string
}

type Redis struct {
	Address  string
	Password string
	Db       int
}

type MongoDB struct {
	WithTimeout int
	DbName      string
	URI         string
}

type Nats struct {
	Url      string
	Username string
	Password string
}

func LoadConfig(fpath string) (*Config, error) {
	if len(fpath) == 0 {
		fpath = "./config.toml"
	}
	var config Config
	_, err := toml.DecodeFile(fpath, &config)
	return &config, err
}
