package src

import (
	"context"
	"fmt"
	"sync"

	// "gorm.io/driver/mysql"
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

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	  if err != nil {
		panic("failed to connect database")
	  }

	return &dbConn{Ctx: ctx, Db: db}
}

func GetDbConn() *dbConn {
	if dbConnInstance == nil {
		lockDbConn.Lock()
		defer lockDbConn.Unlock()
		if dbConnInstance == nil {
			fmt.Println("Creating single instance now.")
			dbConnInstance = createDbConn() 
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return dbConnInstance
}
