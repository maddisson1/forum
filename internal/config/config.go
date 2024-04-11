package config

import "flag"

type Config struct {
	StoragePath string
	Address     string
}

func Loader() *Config {
	addr := flag.String("addr", ":8081", "HTTP network address")
	dsn := flag.String("dsn", "./storage/storage.db", "Sql database storage")
	flag.Parse()

	conf := Config{
		StoragePath: *dsn,
		Address:     *addr,
	}

	return &conf
}
