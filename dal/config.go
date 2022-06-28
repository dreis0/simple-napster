package dal

type Config struct {
	Port     int
	Hostname string
	Username string
	Password string
	Database string
}

func NewConfig(port int, hostname string, username string, password string) *Config {
	return &Config{
		Hostname: hostname,
		Password: password,
		Port:     port,
		Username: username,
	}
}
