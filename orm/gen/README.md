# orm/gen
## 简述
`orm/gen`是一个基于`gorm`,`gorm-gen`的`DB`代码生成工具，它可以根据数据库表结构自动生成`model`,`dao`,`repo`代码，帮助简化数据库操作的开发工作。

## 功能特点
- 自动生成符合gorm规范的数据表结构代码
- 支持多种数据库类型，包括MySQL、PostgreSQL等
- 可自定义数据表命名规则和字段映射规则

## 使用方法
```go
	client, err := orm.NewGormPostgresClient(&orm.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=gorm_gen sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	})
	if err != nil {
		return
	}
	Generation(client, DefaultMySQLDataMap, "./example/postgres/")
```
## 常见问题
Q: 为什么我在运行orm-gen时遇到连接数据库失败的错误？
A: 请确保你已正确配置了数据库连接参数，并确保数据库服务正在运行。同时，检查数据库驱动是否正确安装。

Q: 是否支持其他数据库类型？
A: gorm本身支持多种数据库类型，在orm-gen中理论上同样支持gorm所支持的数据库。如果需要支持其他数据库类型，你可以参考gorm的文档自行扩展。