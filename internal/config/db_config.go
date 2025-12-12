package config

import "time"

type DB struct {
	Driver          string        `required:"true" default:"postgres"`
	Host            string        `required:"true"`
	Port            string        `required:"true"`
	User            string        `required:"true"`
	Password        string        `required:"true"`
	Name            string        `required:"true" default:"cocoon"`
	MaxOpenConns    int           `split_words:"true" default:"25"`
	MaxIdleConns    int           `split_words:"true" default:"5"`
	ConnMaxLifetime time.Duration `split_words:"true" default:"5m"`
}
