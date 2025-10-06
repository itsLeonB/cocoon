package config

type Valkey struct {
	Addr     string `required:"true"`
	Password string `required:"true"`
	Db       int
}
