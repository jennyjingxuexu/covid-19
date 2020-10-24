package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Db DB
}

type DB struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func LoadEnv() (e Env) {
	if err := envconfig.Process("", &e); err != nil {
		panic("Cannot Start App - Env Loading failed")
	}
	return e
}
