package config

type Mail struct {
	Sender   string `required:"true"`
	Host     string `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
	Port     int    `required:"true"`
}
