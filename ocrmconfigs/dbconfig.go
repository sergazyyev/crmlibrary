package ocrmconfigs

import "fmt"

type PostgresDbConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	DbName       string `toml:"db_name"`
	SSLMode      string `toml:"db_ssl_mode"`
	Username     string `toml:"db_username"`
	Password     string `toml:"db_password"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxOpenConns int    `tom:"max_open_conns"`
}

func (p *PostgresDbConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		p.Host,
		p.Port,
		p.DbName,
		p.SSLMode,
		p.Username,
		p.Password)
}
