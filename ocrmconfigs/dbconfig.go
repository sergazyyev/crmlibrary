package ocrmconfigs

type DbConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	DbName       string `toml:"db_name"`
	SSLMode      string `toml:"db_ssl_mode"`
	Username     string `toml:"db_username"`
	Password     string `toml:"db_password"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxOpenConns int    `tom:"max_open_conns"`
}
