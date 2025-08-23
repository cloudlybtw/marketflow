package config

import "os"

type (
	Container struct {
		App       *App
		DB        *DB
		Redis     *Redis
		Exchanges []Exchange
	}

	App struct {
		Name string
		Env  string
	}

	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		Db       string
	}

	Exchange struct {
		Name string
		Port string
	}
)

func New() *Container {
	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Name:       os.Getenv("DB_NAME"),
	}

	redis := &Redis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Db:       os.Getenv("REDIS_DB"),
	}

	exchanges := LoadExchanges()
	return &Container{
		App:       app,
		DB:        db,
		Redis:     redis,
		Exchanges: exchanges,
	}
}

func LoadExchanges() []Exchange {
	return []Exchange{
		loadExchange("EXCHANGE1"),
		loadExchange("EXCHANGE2"),
		loadExchange("EXCHANGE3"),
	}
}

func loadExchange(prefix string) Exchange {
	return Exchange{
		Name: os.Getenv(prefix + "_NAME"),
		Port: os.Getenv(prefix + "_PORT"),
	}
}
