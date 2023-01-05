package utils

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB // 单例

func GetDB() *gorm.DB {
	if db != nil {
		return db
	}

	username := CfgGet("DB_USERNAME")
	password := CfgGet("DB_PASSWORD")
	host := CfgGet("DB_HOST")
	port, _ := strconv.Atoi(CfgGet("DB_PORT"))
	dbname := CfgGet("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbname,
	)

	var err error
	db, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		panic("连接数据库失败，错误信息：" + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("获取底层sqlDB时异常：" + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)           // 空闲连接池的最大数量
	sqlDB.SetMaxOpenConns(100)          // 打开连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接可复用的最大时间

	return db
}
