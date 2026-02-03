package database

import (
	"api_kino/config/app"
	"api_kino/service/env"
	"database/sql/driver"
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/ssh"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Connection string
	Host       string
	Port       int
	Database   string
	Username   string
	Password   string
	IsSSH      bool
}

type ViaSSHDialer struct {
	client *ssh.Client
}

var DB *gorm.DB

func DBConfig() Config {
	return Config{
		Connection: env.GetEnv("DB_CONNECTION"),
		Host:       env.GetEnv("DB_HOST"),
		Port:       env.GetEnvInt("DB_PORT", 5432),
		Database:   env.GetEnv("DB_DATABASE"),
		Username:   env.GetEnv("DB_USERNAME"),
		Password:   env.GetEnv("DB_PASSWORD"),
		IsSSH:      env.GetEnvBool("DB_VIA_SSH", false),
	}
}

func (self *ViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(self, s)
}

func (self *ViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func (self *ViaSSHDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func ConnectAll() {
	DB, _ = Connect(DBConfig())
}

func Connect(dbConfig Config) (db *gorm.DB, er error) {
	dbHost := dbConfig.Host
	dbConnection := dbConfig.Connection
	dbDatabase := dbConfig.Database
	dbUsername := dbConfig.Username
	dbPassword := dbConfig.Password
	logType := logger.Error
	if app.Config().GinMode == gin.DebugMode {
		//logType = logger.Warn
		logType = logger.Info
	}
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logType),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	var connection *gorm.DB
	var err error

	connString := fmt.Sprintf("%s://%s:%s@%s?database=%s", dbConnection, dbUsername, dbPassword, dbHost, dbDatabase)
	connection, err = gorm.Open(sqlserver.Open(connString), gormConfig)

	if err == nil {
		sqlDB, _ := connection.DB()
		sqlDB.SetMaxOpenConns(app.Config().DBMaxOpenCon)
		sqlDB.SetMaxIdleConns(app.Config().DBMaxIdleCon)
		//sqlDB.SetConnMaxLifetime(15 * time.Minute)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	}
	//connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d application_name=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUsername, dbPassword, dbDatabase, dbPort, "api_kino")
	//connection, err := gorm.Open(postgres.New(postgres.Config{
	//	DSN: connString,
	//}), gormConfig)
	//if err == nil {
	//	sqlDB, _ := connection.DB()
	//	sqlDB.SetMaxOpenConns(app.Config().DBMaxOpenCon)
	//	sqlDB.SetMaxIdleConns(app.Config().DBMaxIdleCon)
	//	//sqlDB.SetConnMaxLifetime(15 * time.Minute)
	//	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	//}
	//DB = connection
	return connection, err
}
