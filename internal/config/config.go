package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	App      App      `yaml:"app"`
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"psql"`
	S3       S3       `yaml:"s3"`
	SMTP     SMTP     `yaml:"smtp"`
}

type App struct {
	Cost        int64  `yaml:"cost"`
	FrontendUrl string `yaml:"frontend_url"`
}

type Server struct {
	Scheme      string `yaml:"scheme"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Domain      string `yaml:"domain"`
	SSLCertPath string `yaml:"ssl_cert_path"`
	SSLKeyPath  string `yaml:"ssl_key_path"`
}

type Postgres struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"pass"`
	DbName       string `yaml:"dbname"`
	SSLMode      string `yaml:"sslmode"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type S3 struct {
	Host           string `yaml:"host"`
	AccessKey      string `yaml:"access_key"`
	SecretKey      string `yaml:"secret_key"`
	BucketName     string `yaml:"bucket_name"`
	ServerLocation string `yaml:"server_location"`
	UseSSL         bool   `yaml:"use_ssl"`
}

type SMTP struct {
	Address  string `yaml:"address"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

var config *Config

func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

func Init() (*Config, error) {
	filePath := flag.String("c", "etc/config.yml", "Path to configuration file")
	flag.Parse()
	config = &Config{}
	data, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
