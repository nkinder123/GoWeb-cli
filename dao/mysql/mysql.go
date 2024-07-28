package mysql

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"))
	// 也可以使用MustConnect连接不成功就panic
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.Get("")
	if err != nil {
		zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
		if db != nil {
			Close()
		}
		return
	}
	fmt.Println()
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(viper.GetInt("max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("max_open_conns"))
	return
}

func Close() {
	//由于gorm更新之后不能直接使用.Close()方法，只能将gorm.DB转化成为sql.DB
	var sql *sql.DB
	sql, _ = db.DB()
	sql.Close()

}
func GetDB() *gorm.DB {
	return db
}

//func init() {
//	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
//	dsn := "root:12345678@tcp(localhost:3306)/golangdata?charset=utf8mb4&parseTime=True&loc=Local"
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//}
