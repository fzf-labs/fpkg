package proto

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/fzf-labs/fpkg/orm/gen/utils/file"
	"github.com/fzf-labs/fpkg/orm/gen/utils/template"
	"github.com/fzf-labs/fpkg/orm/gen/utils/util"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GenerationPB 生成
func GenerationPB(db *gorm.DB, outPutPath, packageStr, goPackageStr, table string, columnNameToName map[string]string) error {
	var f string
	p := &Proto{
		gorm:                db,
		outPutPath:          outPutPath,
		packageStr:          packageStr,
		goPackageStr:        goPackageStr,
		tableName:           table,
		tableNameComment:    "",
		tableNameUnderScore: strcase.ToSnake(table),
		lowerTableName:      "",
		upperTableName:      "",
		columnNameToName:    columnNameToName,
	}
	p.tableNameComment = p.GetTableComment(table)
	p.lowerTableName = p.LowerName(table)
	p.upperTableName = p.UpperName(table)
	f += p.GenSyntax()
	f += p.GenPackage()
	f += p.GenImport()
	f += p.GenOption()
	f += p.GenService()
	f += p.GenMessage()
	outputFile := p.outPutPath + "/" + table + ".proto"
	return p.output(outputFile, f)
}

type Proto struct {
	gorm                *gorm.DB          // 数据库
	outPutPath          string            // 生成文件路径
	packageStr          string            // proto中的package名称
	goPackageStr        string            // proto中的goPackage名称
	tableName           string            // 表名称
	tableNameComment    string            // 表注释
	tableNameUnderScore string            // 表下划线名称
	lowerTableName      string            // 表名称首字母小写
	upperTableName      string            // 表名称首字母大写
	columnNameToName    map[string]string // 字段名称对应的Go名称

}

func (p *Proto) output(filePath, content string) error {
	if file.FileExists(filePath) {
		return errors.New(fmt.Sprintf("%s exist", filePath))
	}
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0775); err != nil {
		return err
	}
	dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.WriteString(content)
	if err != nil {
		return err
	}
	return err
}

func (p *Proto) GetTableComment(table string) string {
	type result struct {
		Comment string `json:"comment"`
	}
	var res result
	sql := fmt.Sprintf(`SELECT obj_description(relfilenode,'pg_class')AS comment FROM pg_class WHERE relname='%s'`, table)
	p.gorm.Raw(sql).Scan(&res)
	return res.Comment
}

func (p *Proto) GenSyntax() string {
	str, _ := template.NewTemplate("Syntax").Parse(Syntax).Execute(map[string]any{})
	return fmt.Sprintln(str.String())
}

func (p *Proto) GenPackage() string {
	str, _ := template.NewTemplate("Package").Parse(Package).Execute(map[string]any{
		"packageStr": p.packageStr,
	})
	return fmt.Sprintln(str.String())
}

func (p *Proto) GenImport() string {
	str, _ := template.NewTemplate("Import").Parse(Import).Execute(map[string]any{})
	return fmt.Sprintln(str.String())
}

func (p *Proto) GenOption() string {
	str, _ := template.NewTemplate("Option").Parse(Option).Execute(map[string]any{
		"goPackageStr": p.goPackageStr,
	})
	return fmt.Sprintln(str.String())
}

func (p *Proto) GenService() string {
	str, _ := template.NewTemplate("Service").Parse(Service).Execute(map[string]any{
		"upperTableName":      p.upperTableName,
		"tableNameComment":    p.tableNameComment,
		"tableNameUnderScore": p.tableNameUnderScore,
	})
	return fmt.Sprintln(str.String())
}

func (p *Proto) GenMessage() string {
	var info string
	var createReq string
	var createReply string
	var updateReq string
	var deleteReq string
	var getReq string
	columnTypes, err := p.gorm.Migrator().ColumnTypes(p.tableName)
	if err != nil {
		return ""
	}
	indexes, err := p.gorm.Migrator().GetIndexes(p.tableName)
	if err != nil {
		return ""
	}
	var primaryKeyColumn string
	for _, index := range indexes {
		b, _ := index.PrimaryKey()
		if b {
			primaryKeyColumn = index.Columns()[0]
			break
		}
	}
	columnTypeInfo := make(map[string]gorm.ColumnType)
	num := 0
	createNum := 0
	for _, v := range columnTypes {
		num++
		columnTypeInfo[v.Name()] = v
		pbType := columnTypeToPbType(v.DatabaseTypeName())
		pbName := LowerFieldName(p.columnNameToName[v.Name()])
		comment, _ := v.Comment()
		if util.StrSliFind([]string{"deletedAt", "deleted_at", "deletedTime", "deleted_time"}, v.Name()) {
			continue
		}
		info += fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, num, comment)
		if util.StrSliFind([]string{"createdAt", "created_at", "createdTime", "created_time", "updatedAt", "updated_at", "updatedTime", "updated_time"}, v.Name()) {
			continue
		}
		if v.Name() != primaryKeyColumn {
			createNum++
			createReq += fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, createNum, comment)
		}
		updateReq += fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, num, comment)
	}
	if primaryKeyColumn != "" {
		primaryKeyColumnType, _ := columnTypeInfo[primaryKeyColumn].ColumnType()
		primaryKeyComment, _ := columnTypeInfo[primaryKeyColumn].Comment()
		pbType := columnTypeToPbType(primaryKeyColumnType)
		pbName := LowerFieldName(p.columnNameToName[primaryKeyColumn])
		createReply = fmt.Sprintf("	%s %s = %d; // %s", pbType, pbName, 1, primaryKeyComment)
		getReq = fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, 1, primaryKeyComment)
		deleteReq = fmt.Sprintf("repeated %s %s = %d; // %s\n", pbType, Plural(pbName), 1, primaryKeyComment+"集合")
	}
	info = strings.TrimSpace(strings.TrimRight(info, "\n"))
	createReq = strings.TrimSpace(strings.TrimRight(createReq, "\n"))
	updateReq = strings.TrimSpace(strings.TrimRight(updateReq, "\n"))
	deleteReq = strings.TrimSpace(strings.TrimRight(deleteReq, "\n"))
	getReq = strings.TrimSpace(strings.TrimRight(getReq, "\n"))
	str, _ := template.NewTemplate("Message").Parse(Message).Execute(map[string]any{
		"tableNameComment": p.tableNameComment,
		"upperTableName":   p.upperTableName,
		"info":             info,
		"createReq":        createReq,
		"createReply":      createReply,
		"updateReq":        updateReq,
		"deleteReq":        deleteReq,
		"getReq":           getReq,
	})
	return fmt.Sprintln(str.String())
}

// UpperName 大写
func (p *Proto) UpperName(s string) string {
	return p.gorm.NamingStrategy.SchemaName(s)
}

// LowerName 小写
func (p *Proto) LowerName(s string) string {
	str := p.UpperName(s)
	if str == "" {
		return str
	}
	words := []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "ttl", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	// 如果第一个单词命中  则不处理
	for _, v := range words {
		if strings.HasPrefix(str, v) {
			return str
		}
	}
	rs := []rune(str)
	f := rs[0]
	if 'A' <= f && f <= 'Z' {
		str = string(unicode.ToLower(f)) + string(rs[1:])
	}
	return str
}

// LowerFieldName 字段名称小写
func LowerFieldName(str string) string {
	words := []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "ttl", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	// 如果第一个单词命中  则不处理
	for _, v := range words {
		if strings.HasPrefix(str, v) {
			return str
		}
	}
	rs := []rune(str)
	f := rs[0]
	if 'A' <= f && f <= 'Z' {
		str = string(unicode.ToLower(f)) + string(rs[1:])
	}
	if token.Lookup(str).IsKeyword() {
		str = "_" + str
	}
	return str
}

func columnTypeToPbType(columnType string) string {
	var fieldType string
	switch columnType {
	case "char", "varchar", "text", "uuid", "json", "jsonb":
		fieldType = "string"
	case "date", "timestamp", "timetz", "timestamptz":
		fieldType = "google.protobuf.Timestamp"
	case "bool":
		fieldType = "bool"
	case "int2", "int4", "int8":
		fieldType = "int32"
	case "float4":
		fieldType = "float"
	case "float8":
		fieldType = "double"
	default:
		fieldType = "string"
	}
	return fieldType
}

// Plural 复数形式
func Plural(s string) string {
	str := inflection.Plural(s)
	if str == s {
		str += "Plural"
	}
	return str
}
