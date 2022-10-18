package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	DB *gorm.DB
}

func NewGormDB() *GormDB {
	return &GormDB{}
}

func (g *GormDB) Open(uri string, cfg *gorm.Config) error {
	db, err := gorm.Open(
		postgres.Open(uri),
		cfg,
	)
	if err != nil {
		return err
	}

	g.DB = db

	return nil
}

func (g *GormDB) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}
