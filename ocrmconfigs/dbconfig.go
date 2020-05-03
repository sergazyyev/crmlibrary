package ocrmconfigs

import (
	"fmt"
	"net/url"
	"time"
)

type PostgresDbConfig struct {
	Host            string        `toml:"host"`
	Port            int           `toml:"port"`
	DbName          string        `toml:"db_name"`
	SSLMode         string        `toml:"db_ssl_mode"`
	Username        string        `toml:"db_username"`
	Password        string        `toml:"db_password"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
	MaxOpenConns    int           `toml:"max_open_conns"`
	ConnMaxLifeTime time.Duration `toml:"conn_max_life_time"`
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

type OracleDbConfig struct {
	Host            string        `toml:"host"`
	Port            int           `toml:"port"`
	DbName          string        `toml:"db_name"`
	Username        string        `toml:"db_username"`
	Password        string        `toml:"db_password"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
	MaxOpenConns    int           `toml:"max_open_conns"`
	ConnMaxLifeTime time.Duration `toml:"conn_max_life_time"`
}

func (o *OracleDbConfig) GetConnectionString() string {
	return fmt.Sprintf(`oracle://%s:%s@%s:%d/%s`,
		url.QueryEscape(o.Username),
		url.QueryEscape(o.Password),
		url.QueryEscape(o.Host),
		o.Port,
		url.QueryEscape(o.DbName))
}
