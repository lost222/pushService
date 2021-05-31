package model

import (
"fmt"
_ "github.com/go-sql-driver/mysql"
"github.com/jinzhu/gorm"
)

var db *gorm.DB

var err error


var (
	DbUser string = "root"
	DbPassWord string = "root"
	DbHost string = "localhost"
	DbPort string = "3306"
	DbName = "ginrss"
	Db = "mysql"
)
func InitDB()  {
	dmm := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DbUser,
		DbPassWord,
		DbHost,
		DbPort,
		DbName,
	)
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(Db, dmm)

	if err != nil{
		fmt.Println("connect to mysql wrong", err)
	}

}



func LimString(maxLen int, s string) string{
	if len(s) < maxLen{
		return s
	}
	return s[:maxLen]
}


