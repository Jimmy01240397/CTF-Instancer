package database

import (
    "log"
    "time"
    "fmt"

    "gorm.io/driver/sqlite"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm/logger"
    "gorm.io/gorm"
    
    "github.com/Jimmy01240397/CTF-Instancer/utils/config"
)

var db *gorm.DB

func init() {
    var err error
    gormconfig := &gorm.Config{}
    if config.Debug {
        gormconfig.Logger = logger.Default.LogMode(logger.Info)
    }
    if !config.DBdebug {
        gormconfig.Logger = logger.Default.LogMode(logger.Silent)
    }
    switch config.DBservice {
    case "sqlite":
        db, err = gorm.Open(sqlite.Open(config.DBname), gormconfig)
    case "mysql":
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config.DBuser, config.DBpasswd, config.DBhost, config.DBport, config.DBname)
        db, err = gorm.Open(mysql.Open(dsn), gormconfig)
    case "postgres":
        dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei", config.DBhost, config.DBuser, config.DBpasswd, config.DBname, config.DBport)
        db, err = gorm.Open(postgres.Open(dsn), gormconfig)
    }
    if err != nil {
        log.Panicln(err)
    }
    sqlDB, err := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetDB() *gorm.DB {
    return db
}
