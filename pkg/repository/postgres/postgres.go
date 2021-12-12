package postgres

import (
	"fmt"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	Database string `mapstructure:"database"`
	Port     string `mapstructure:"port"`
}

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(cfg *Config) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s sslmode=disable", cfg.Address, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	p := &Postgres{db: db}
	return p, nil
}

func (p *Postgres) Migrate() error {
	err := p.db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}
