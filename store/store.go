package store

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/teacat/goshia/v3"
	"github.com/teacatx/iknore/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Store struct {
	Goshia *goshia.Goshia
}

type Conn struct {
	DB *gorm.DB
}

func NewConn() *Conn {
	if v := os.Getenv("DB_USERNAME"); v != "" {
		config.C.DBUsername = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		config.C.DBPassword = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		config.C.DBHost = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		config.C.DBPort = v
	}
	if v := os.Getenv("DB_DATABASE"); v != "" {
		config.C.DBDatabase = v
	}
	if v := os.Getenv("DB_CHARSET"); v != "" {
		config.C.DBCharset = v
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.C.DBUsername,
		config.C.DBPassword,
		config.C.DBHost,
		config.C.DBPort,
		config.C.DBDatabase,
		config.C.DBCharset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      false,
			},
		),
	})
	if err != nil {
		log.Fatal(err)
	}
	//
	db.AutoMigrate(
		&Image{},
		&Pointer{},
	)
	//
	return &Conn{
		DB: db,
	}
}

// New
func New(g *goshia.Goshia) *Store {
	return &Store{
		Goshia: g,
	}
}
