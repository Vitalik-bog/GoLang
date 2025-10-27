package config

type Config struct {
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL: "host=localhost port=5433 user=postgres password=********** dbname=employees sslmode=disable",
	}
}
