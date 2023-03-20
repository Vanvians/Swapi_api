package db

import (
    "fmt"
    "log"

    "github.com/go-redis/redis"
    "github.com/jcezetah/Swapi_api/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (
    db     *gorm.DB
    rdb    *redis.Client
    dbHost = "localhost"
    dbPort = 5432
    dbName = "swapi_db"
    dbUser = "postgres"
    dbPass = "password"
)

func InitDB() (*gorm.DB, *redis.Client, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPass, dbName)

    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database", err)
    }

    // Auto migrate the models
    db.AutoMigrate(&models.Movie{}, &models.Comment{})

    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    _, err = rdb.Ping().Result()
    if err != nil {
        log.Fatal("failed to connect redis", err)
    }

    return db, rdb, nil
}
