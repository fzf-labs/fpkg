package orm

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fzf-labs/fpkg/orm/plugin"
	"github.com/fzf-labs/fpkg/util/fileutil"
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

// DumpPostgres 导出创建语句
func DumpPostgres(dsn string, tables []string, outPath string) {
	// 查找命令的可执行文件
	_, err := exec.LookPath("pg_dump")
	if err != nil {
		slog.Error("command pg_dump not found,please install")
		return
	}
	if len(tables) == 0 {
		db, err2 := gorm.Open(postgres.Open(dsn))
		if err2 != nil {
			log.Fatal("open db err2:", err2.Error())
			return
		}
		tables, err2 = db.Migrator().GetTables()
		if err2 != nil {
			return
		}
	}
	dsnParse := PostgresDsnParse(dsn)
	outPath = filepath.Join(strings.Trim(outPath, "/"), dsnParse.Dbname)
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		slog.Error("DumpPostgres create path err:", err)
		return
	}
	for _, v := range tables {
		outFile := filepath.Join(outPath, fmt.Sprintf("%s.sql", v))
		cmdArgs := []string{
			"-h", dsnParse.Host,
			"-p", strconv.Itoa(dsnParse.Port),
			"-U", dsnParse.User,
			"-s", dsnParse.Dbname,
			"-t", v,
		}
		// 创建一个 Cmd 对象来表示将要执行的命令
		cmd := exec.Command("pg_dump", cmdArgs...)
		// 添加一个环境变量到命令中
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", dsnParse.Password))
		// 执行命令，并捕获输出和错误信息
		output, err := cmd.Output()
		if err != nil {
			slog.Error("cmd exec err:", err)
			return
		}
		err = fileutil.WriteContentCover(outFile, remove(string(output)))
		if err != nil {
			slog.Error("DumpPostgres err:", err)
			return
		}
	}
}

type PostgresDsn struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

// PostgresDsnParse  数据库解析
func PostgresDsnParse(dsn string) *PostgresDsn {
	result := new(PostgresDsn)
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
			result.Host = value
		case "port":
			if p, err := strconv.Atoi(value); err == nil {
				result.Port = p
			}
		case "user":
			result.User = value
		case "password":
			result.Password = value
		case "dbname":
			result.Dbname = value
		}
	}
	return result
}

// remove 移除多余行
func remove(str string) string {
	var result string
	reader := strings.NewReader(str)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "SELECT") || strings.HasPrefix(line, "SET") || regexp.MustCompile(`(ALTER TABLE .*? OWNER TO postgres)`).MatchString(line) {
			continue
		}
		result += fmt.Sprintln(line)
	}
	return result
}
