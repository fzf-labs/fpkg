package pb

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
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Pb struct {
	gorm         *gorm.DB
	pbPath       string
	packageStr   string
	goPackageStr string
}

func NewPbRepo(db *gorm.DB, pbPath, packageStr, goPackageStr string) *Pb {
	return &Pb{
		gorm:         db,
		pbPath:       pbPath,
		packageStr:   packageStr,
		goPackageStr: goPackageStr,
	}
}

// GenerationTable 生成
func (r *Pb) GenerationTable(table string, columnNameToName map[string]string) error {
	var file string
	g := &GenerationPb{
		gorm:                r.gorm,
		tableName:           table,
		tableNameComment:    "",
		tableNameUnderScore: strcase.ToSnake(table),
		lowerTableName:      "",
		upperTableName:      "",
		columnNameToName:    columnNameToName,
		packageStr:          r.packageStr,
		goPackageStr:        r.goPackageStr,
	}
	g.tableNameComment = g.GetTableComment(table)
	g.lowerTableName = g.LowerName(table)
	g.upperTableName = g.UpperName(table)
	file += g.GenSyntax()
	file += g.GenPackage()
	file += g.GenImport()
	file += g.GenOption()
	file += g.GenService()
	file += g.GenMessage()
	outputFile := r.pbPath + "/" + table + ".proto"
	return r.output(outputFile, file)
}

func (r *Pb) output(filePath, content string) error {
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

type GenerationPb struct {
	gorm                *gorm.DB
	tableName           string            // 表名称
	tableNameComment    string            // 表注释
	tableNameUnderScore string            // 表下划线名称
	lowerTableName      string            // 表名称首字母小写
	upperTableName      string            // 表名称首字母大写
	columnNameToName    map[string]string // 字段名称对应的Go名称
	packageStr          string
	goPackageStr        string
}

func (g *GenerationPb) GetTableComment(table string) string {
	var name string
	sql := fmt.Sprintf(`SELECT obj_description(relfilenode,'pg_class')AS table_comment FROM pg_class WHERE relname='%s'`, table)
	g.gorm.Raw(sql).Scan(&name)
	return name
}

func (g *GenerationPb) GenSyntax() string {
	str, _ := template.NewTemplate("Syntax").Parse(Syntax).Execute(map[string]any{})
	return fmt.Sprintln(str.String())
}

func (g *GenerationPb) GenPackage() string {
	str, _ := template.NewTemplate("Package").Parse(Package).Execute(map[string]any{
		"packageStr": g.packageStr,
	})
	return fmt.Sprintln(str.String())
}

func (g *GenerationPb) GenImport() string {
	str, _ := template.NewTemplate("Import").Parse(Import).Execute(map[string]any{})
	return fmt.Sprintln(str.String())
}

func (g *GenerationPb) GenOption() string {
	str, _ := template.NewTemplate("Option").Parse(Option).Execute(map[string]any{
		"goPackageStr": g.goPackageStr,
	})
	return fmt.Sprintln(str.String())
}

func (g *GenerationPb) GenService() string {
	str, _ := template.NewTemplate("Service").Parse(Service).Execute(map[string]any{
		"upperTableName":      g.upperTableName,
		"tableNameComment":    g.tableNameComment,
		"tableNameUnderScore": g.tableNameUnderScore,
	})
	return fmt.Sprintln(str.String())
}

func (g *GenerationPb) GenMessage() string {
	columnTypes, err := g.gorm.Migrator().ColumnTypes(g.tableName)
	if err != nil {
		return ""
	}
	var info string
	var storeReq string
	var storeReply string
	var delReq string
	var oneReq string
	columnTypeInfo := make(map[string]gorm.ColumnType)
	for k, v := range columnTypes {
		columnTypeInfo[v.Name()] = v
		pbType := columnTypeToPbType(v.DatabaseTypeName())
		pbName := LowerFieldName(g.columnNameToName[v.Name()])
		comment, _ := v.Comment()
		if util.StrSliFind([]string{"deletedAt", "deleted_at", "deletedTime", "deleted_time"}, v.Name()) {
			continue
		}
		info += fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, k+1, comment)
		if util.StrSliFind([]string{"createdAt", "created_at", "createdTime", "created_time", "updatedAt", "updated_at", "updatedTime", "updated_time"}, v.Name()) {
			continue
		}
		storeReq += fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, k+1, comment)
	}
	indexes, err := g.gorm.Migrator().GetIndexes(g.tableName)
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
	if primaryKeyColumn != "" {
		primaryKeyColumnType, _ := columnTypeInfo[primaryKeyColumn].ColumnType()
		primaryKeyComment, _ := columnTypeInfo[primaryKeyColumn].Comment()
		pbType := columnTypeToPbType(primaryKeyColumnType)
		pbName := LowerFieldName(g.columnNameToName[primaryKeyColumn])
		storeReply = fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, 1, primaryKeyComment)
		oneReq = fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, 1, primaryKeyComment)
		delReq = fmt.Sprintf("	%s %s = %d; // %s\n", pbType, pbName, 1, primaryKeyComment)
	}
	str, _ := template.NewTemplate("Message").Parse(Message).Execute(map[string]any{
		"tableNameComment": g.tableNameComment,
		"upperTableName":   g.upperTableName,
		"Info":             info,
		"StoreReq":         storeReq,
		"StoreReply":       storeReply,
		"DelReq":           delReq,
		"OneReq":           oneReq,
	})
	return fmt.Sprintln(str.String())
}

// UpperName 大写
func (g *GenerationPb) UpperName(s string) string {
	return g.gorm.NamingStrategy.SchemaName(s)
}

// LowerName 小写
func (g *GenerationPb) LowerName(s string) string {
	str := g.UpperName(s)
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
