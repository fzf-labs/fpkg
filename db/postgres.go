package db

import (
	"database/sql"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/fzf-labs/fpkg/db/plugin"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

// GormPostgresClientConfig 配置
type GormPostgresClientConfig struct {
	DataSourceName  string        `json:"DataSourceName"`
	MaxIdleConn     int           `json:"MaxIdleConn"`
	MaxOpenConn     int           `json:"MaxOpenConn"`
	ConnMaxLifeTime time.Duration `json:"ConnMaxLifeTime"`
	ShowLog         bool          `json:"ShowLog"`
	Tracing         bool          `json:"Tracing"`
	Caches          bool          `json:"Caches"`
}

// NewGormPostgresClient 初始化gorm Postgres 客户端
func NewGormPostgresClient(cfg *GormPostgresClientConfig) (*gorm.DB, error) {
	sqlDB, err := sql.Open("pgx", cfg.DataSourceName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("open mysql failed! err: %+v", err))
	}
	// set for db connection
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	// 设置连接可以重复使用的最长时间.
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}
	if cfg.ShowLog {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gormConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("postgres database connection failed!  err: %+v", err))
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	if cfg.Tracing {
		if err := db.Use(tracing.NewPlugin()); err != nil {
			return nil, errors.New(fmt.Sprintf("db use tracing failed!  err: %+v", err))
		}
	}
	if cfg.Caches {
		if err := db.Use(plugin.NewCaches()); err != nil {
			return nil, errors.New(fmt.Sprintf("db use Caches failed!  err: %+v", err))
		}
	}
	return db, nil
}

// DumpSql 导出创建语句
func DumpSql(db *gorm.DB, dsn string, outPath string) {
	// 查找命令的可执行文件
	_, err := exec.LookPath("pg_dump")
	if err != nil {
		fmt.Printf("Command %s not found\n", "pg_dump")
		return
	}
	tables, err := db.Migrator().GetTables()
	if err != nil {
		return
	}
	outPath = strings.Trim(outPath, "/")
	dsnParse := DsnParse(dsn)
	for _, v := range tables {
		cmdArgs := []string{
			"-h", dsnParse.host,
			"-p", strconv.Itoa(dsnParse.port),
			"-U", dsnParse.user,
			"-s", dsnParse.dbname,
			"-t", v,
			"-f", fmt.Sprintf("%s/%s.sql", outPath, v),
		}
		// 创建一个 Cmd 对象来表示将要执行的命令
		cmd := exec.Command("pg_dump", cmdArgs...)
		// 添加一个环境变量到命令中
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", dsnParse.password))
		// 执行命令，并捕获输出和错误信息
		err := cmd.Run()
		if err != nil {
			fmt.Println("DumpSql err:", err)
			return
		}
	}
}

type Dsn struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func DsnParse(dsn string) *Dsn {
	result := new(Dsn)
	// 分割连接字符串
	params := strings.Split(dsn, " ")

	// 解析参数
	for _, param := range params {
		keyValue := strings.Split(param, "=")
		if len(keyValue) != 2 {
			continue
		}
		key := keyValue[0]
		value := keyValue[1]
		switch key {
		case "host":
			result.host = value
		case "port":
			if p, err := strconv.Atoi(value); err == nil {
				result.port = p
			}
		case "user":
			result.user = value
		case "password":
			result.password = value
		case "dbname":
			result.dbname = value
		}
	}
	return result
}
