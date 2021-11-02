package datasource

import (
	"fmt"
	"log"

	"github.com/gtck520/ConsoleIM/common/logger"
	"github.com/gtck520/ConsoleIM/common/setting"
	"github.com/gtck520/ConsoleIM/models/common"
	"github.com/gtck520/ConsoleIM/models/db"
	"github.com/gtck520/ConsoleIM/models/sys"
	"github.com/jinzhu/gorm"

	//
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Db gormDB
type Db struct {
	Conn *gorm.DB
}

//Connect 初始化数据库配置
func (d *Db) Connect() error {
	var (
		dbType, dbName, user, pwd, host string
	)

	conf := setting.Database
	dbType = conf["Type"]
	dbName = conf["Name"]
	user = conf["User"]
	pwd = conf["Password"]
	host = conf["Host"] + ":" + conf["Port"]

	gdb, err := gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd, host, dbName))
	if err != nil {
		log.Fatal("connecting mysql error: ", err)
		return err
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		//return conf["Prefix"] + defaultTableName
		return defaultTableName
	}
	gdb.LogMode(true) //打印SQL语句
	gdb.SingularTable(true)
	gdb.SetLogger(logger.SqlLogger())
	gdb.DB().SetMaxIdleConns(10)
	gdb.DB().SetMaxOpenConns(100)

	d.Conn = gdb
	db.DB = gdb

	log.Println("Connect Mysql Success")

	return nil
}

//DB 返回DB
func (d *Db) DB() *gorm.DB {
	return d.Conn
}

//自动生成数据库
func Migration() {
	fmt.Println(db.DB.AutoMigrate(new(sys.Menu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.Admins)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.RoleMenu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.Role)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.AdminsRole)).Error)
	fmt.Println(db.DB.AutoMigrate(new(common.User)).Error)
	fmt.Println(db.DB.AutoMigrate(new(common.Friend)).Error)
	fmt.Println(db.DB.AutoMigrate(new(common.Group)).Error)
	fmt.Println(db.DB.AutoMigrate(new(common.UserGroup)).Error)
}
