package config

import "time"

type DB struct {
	Driver          string        `required:"true" default:"postgres"`
	Host            string        `required:"true"`
	Port            string        `required:"true"`
	User            string        `required:"true"`
	Password        string        `required:"true"`
	Name            string        `required:"true" default:"cocoon"`
	MaxOpenConns    int           `default:"25"`
	MaxIdleConns    int           `default:"5"`
	ConnMaxLifetime time.Duration `default:"5m"`
}
