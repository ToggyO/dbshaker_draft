package suites

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
)

type DbConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewDbConf(files ...string) DbConf {
	var err error

	absPaths := make([]string, 0, cap(files))
	for _, f := range files {
		envPath, err := filepath.Abs(f)
		if err != nil {
			panic(err)
		}
		absPaths = append(absPaths, envPath)
	}

	err = godotenv.Load(absPaths...)
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return DbConf{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
