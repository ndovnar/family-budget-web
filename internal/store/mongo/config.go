package mongo

type Config struct {
	Hosts    []string `required:"true" desc:"mongo hosts"`
	Username string   `required:"true" desc:"username for mongo access"`
	Password string   `required:"true" desc:"password for mongo access" secret:"true"`
	Database string   `required:"true" desc:"mongo database name"`
	UseTLS   bool     `desc:"enable TLS on all database connections"`
}
