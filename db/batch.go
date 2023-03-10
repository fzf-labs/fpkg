package db

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// BatchUpdateToSqlArray 批量更新
func BatchUpdateToSqlArray(tableName string, dataList interface{}) ([]string, error) {
	//dataList 的数据类型
	kind := reflect.TypeOf(dataList).Kind()
	if kind != reflect.Slice {
		return nil, errors.New("please set dataList for sli")
	}
	fieldValue := reflect.ValueOf(dataList)
	fieldType := reflect.TypeOf(dataList).Elem().Elem()
	//切片的长度
	sliceLength := fieldValue.Len()
	//字段的数量
	fieldNum := fieldType.NumField()

	// 检验结构体标签是否为空和重复
	verifyTagDuplicate := make(map[string]string)
	for i := 0; i < fieldNum; i++ {
		//获取json tag字段的值
		fieldName := strings.TrimSpace(fieldType.Field(i).Tag.Get("json"))
		if fieldName == "" {
			return nil, errors.New("please set json tag for struct")
		}
		_, ok := verifyTagDuplicate[fieldName]
		if !ok {
			verifyTagDuplicate[fieldName] = fieldName
		} else {
			return nil, errors.New(fmt.Sprintf("the structure tag: %v  is not allow duplication", fieldName))
		}
	}
	//校验是否有id字段
	if _, ok := verifyTagDuplicate["id"]; !ok {
		return nil, errors.New("please set json tag: id for struct")
	}

	var Ids []string
	updateMap := make(map[string][]string)
	for i := 0; i < sliceLength; i++ {
		// 得到某一个具体的结构体的
		structValue := fieldValue.Index(i).Elem()
		for j := 0; j < fieldNum; j++ {
			elem := structValue.Field(j)
			var temp string
			switch elem.Kind() {
			case reflect.Int64:
				temp = strconv.FormatInt(elem.Int(), 10)
			case reflect.String:
				if strings.Contains(elem.String(), "'") {
					temp = fmt.Sprintf("'%v'", strings.ReplaceAll(elem.String(), "'", "\\'"))
				} else {
					temp = fmt.Sprintf("'%v'", elem.String())
				}
			case reflect.Float64:
				temp = strconv.FormatFloat(elem.Float(), 'f', -1, 64)
			case reflect.Bool:
				temp = strconv.FormatBool(elem.Bool())
			default:
				return nil, errors.New(fmt.Sprintf("type conversion error, param is %v", fieldType.Field(j).Tag.Get("json")))
			}
			fieldName := strings.TrimSpace(fieldType.Field(j).Tag.Get("json"))
			if fieldName == "id" {
				id, err := strconv.ParseInt(temp, 10, 64)
				if err != nil {
					return nil, err
				}
				// id 的合法性校验
				if id < 1 {
					return nil, errors.New("this structure should have a primary key and gt 0")
				}
				Ids = append(Ids, temp)
				continue
			}
			valueList := append(updateMap[fieldName], temp)
			updateMap[fieldName] = valueList
		}
	}

	length := len(Ids)
	size := 200
	SQLQuantity := getSQLQuantity(length, size)
	var SQLArray []string
	k := 0
	for i := 0; i < SQLQuantity; i++ {
		count := 0
		var record bytes.Buffer
		record.WriteString("UPDATE " + tableName + " SET ")
		for fieldName, fieldValueList := range updateMap {
			record.WriteString(fieldName)
			record.WriteString(" = CASE " + "id")
			for j := k; j < len(Ids) && j < len(fieldValueList) && j < size+k; j++ {
				record.WriteString(" WHEN " + Ids[j] + " THEN " + fieldValueList[j])
			}
			count++
			if count != fieldNum-1 {
				record.WriteString(" END, ")
			}
		}
		record.WriteString(" END WHERE ")
		record.WriteString("id" + " IN (")
		min := size + k
		if len(Ids) < min {
			min = len(Ids)
		}
		record.WriteString(strings.Join(Ids[k:min], ","))
		record.WriteString(")")
		k += size
		SQLArray = append(SQLArray, record.String())
	}
	return SQLArray, nil
}

func getSQLQuantity(length, size int) int {
	SQLQuantity := int(math.Ceil(float64(length) / float64(size)))
	return SQLQuantity
}
