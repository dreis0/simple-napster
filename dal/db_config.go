package dal

import "simple-napster/utils"

type DbConfig struct {
	Url      string
	Username string
	Database string
	Password string
}

func FromEnv() *DbConfig {
	dbUrl := utils.FromEnv("DATABASE_URL", "localhost:5432")
	dbUser := utils.FromEnv("DB_USER", "postgres")
	dbName := utils.FromEnv("DB_NAME", "napster")
	dbPassword := utils.FromEnv("DB_PASSWORD", "can't use default")

	return &DbConfig{
		Url:      dbUrl,
		Username: dbUser,
		Database: dbName,
		Password: dbPassword,
	}
}
