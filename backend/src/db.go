package src

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aarhunt/spootify/src/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var lockDbConn = &sync.Mutex{}

type dbConn struct {
	Ctx context.Context
	Db 	*gorm.DB
}

var dbConnInstance *dbConn

func createDbConn() *dbConn{
	ctx := context.Background()

	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	env := os.Getenv("NODE_ENV")
	fmt.Println(env)

	var db *gorm.DB
	var err error
	if env == "production" {
		db, err = gorm.Open(postgres.Open(pgConnString), &gorm.Config{})
		  if err != nil {
			panic("failed to connect database")
		  }
	} else if env == "development" {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		  if err != nil {
			panic("failed to connect database")
		  }
	}

	db.AutoMigrate(&model.Playlist{})

	return &dbConn{Ctx: ctx, Db: db}
}

func GetDbConn() *dbConn {
	if dbConnInstance == nil {
		lockDbConn.Lock()
		defer lockDbConn.Unlock()
		if dbConnInstance == nil {
			fmt.Println("Creating single instance now.")
			dbConnInstance = createDbConn() 
		}
	}

	return dbConnInstance
}
