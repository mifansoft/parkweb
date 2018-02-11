package orm

import (
	"mifanpark/utilities/crypto/aes"
	"mifanpark/utilities/logger"
	"mifanpark/utilities/orm/dbhandle"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {
	cfg, err := dbhandle.GetConfig()
	if err != nil {
		panic("cant not read ./conf/app.conf.please check this file.")
	}
	dbDriver, _ := cfg.Get("DB.Driver")
	dbhandle.Register(dbDriver, NewDB)
}

func NewDB() gorm.DB {
	var err error
	var db *gorm.DB
	cfg, err := dbhandle.GetConfig()
	if err != nil {
		panic("cant not read ./conf/app.conf.please check this file.")
	}
	dbHost, _ := cfg.Get("DB.Host")
	dbUserName, _ := cfg.Get("DB.UserName")
	dbPassword, _ := cfg.Get("DB.Password")
	dbPort, _ := cfg.Get("DB.Port")
	dbDatabase, _ := cfg.Get("DB.Database")
	dbCharset, _ := cfg.Get("DB.Charset")
	dbProtocol, _ := cfg.Get("DB.Protocol")
	dbDriver, _ := cfg.Get("DB.Driver")

	if len(dbPassword) == 24 {
		dbPassword, err = aes.Decrypt(dbPassword)
		if err != nil {
			logger.Error("Decrypt mysql password failed.")
			return *db
		}
	}
	var conn string
	switch dbDriver {
	case "mysql":
		conn = dbUserName + ":" + dbPassword + "@" + dbProtocol + "(" + dbHost + ":" + dbPort + ")/" + dbDatabase + "?charset=" + dbCharset + "&parseTime=true&loc=Local"
	}
	db, err = gorm.Open(dbDriver, conn)
	if err != nil {
		logger.Error(err)
	}
	if len(dbPassword) != 24 {
		psd, err := aes.Encrypt(dbPassword)
		if err != nil {
			logger.Error("descrypt password failed." + dbPassword)
			return *db
		}
		psd = "\"" + dbPassword + "\""
		cfg.Set("DB.Password", psd)
	}
	db.SingularTable(true)
	return *db
}
