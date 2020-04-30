package ocrmconfigs

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/godror/godror"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
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

func (p *PostgresDbConfig) ConfigureDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", p.GetConnectionString())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if p.MaxIdleConns != 0 {
		db.SetMaxIdleConns(p.MaxIdleConns)
	}
	if p.MaxOpenConns != 0 {
		db.SetMaxOpenConns(p.MaxOpenConns)
	}
	if p.ConnMaxLifeTime != 0 {
		db.SetConnMaxLifetime(p.ConnMaxLifeTime)
	}
	return db, nil
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

func (o *OracleDbConfig) ConfigureDb() (*sqlx.DB, error) {
	db, err := sqlx.Open("godror", o.GetConnectionString())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if o.MaxIdleConns != 0 {
		db.SetMaxIdleConns(o.MaxIdleConns)
	}
	if o.MaxOpenConns != 0 {
		db.SetMaxOpenConns(o.MaxOpenConns)
	}
	if o.ConnMaxLifeTime != 0 {
		db.SetConnMaxLifetime(o.ConnMaxLifeTime)
	}
	return db, nil
}
