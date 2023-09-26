package tools

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"strings"
	"time"
)

func GetExcelElement(inputSlice []string) string {
	result := ""
	for i, element := range inputSlice {
		// Split the element by ", " and join them with "\t"
		element = strings.Split(element, ".")[1]
		parts := strings.Split(element, ", ")
		joinedElement := strings.Join(parts, "\t")

		// Append the joined element to the result string
		result += joinedElement

		// Add a newline character to separate rows except for the last element
		if i < len(inputSlice)-1 {
			result += "\r\n"
		}
	}
	return result
}

func StmtToGo(stmt *sqlparser.DDL, tableName string, pkgName string) ([]string, error) {
	builder := strings.Builder{}

	// header := fmt.Sprintf("package %s\n", pkgName)

	structName := snakeCaseToCamel(tableName)
	structStart := fmt.Sprintf("type %s struct { \n", structName)
	builder.WriteString(structStart)
	ret := make([]string, 0)
	for _, col := range stmt.TableSpec.Columns {
		columnType := col.Type.Type

		goType := sqlTypeMap[columnType]

		field := snakeCaseToCamel(col.Name.String())
		retStr := field + ", " + goType
		comment := col.Type.Comment
		if comment == nil {
			builder.WriteString(fmt.Sprintf("\t%s\t%s\t\n", field, goType))
			retStr = "-" + ", " + retStr
		} else {
			builder.WriteString(fmt.Sprintf("\t%s\t%s\t`comment:\"%s\"` \n",
				field, goType, string(comment.Val)))
			retStr = string(comment.Val) + ", " + retStr
		}
		if retStr == "" {
			continue
		}
		ret = append(ret, retStr)
	}
	builder.WriteString("}\n")

	return ret, nil
}

// In sql, table name often is snake_case
// In Go, struct name often is camel case
func snakeCaseToCamel(str string) string {
	builder := strings.Builder{}
	index := 0
	if str[0] >= 'a' && str[0] <= 'z' {
		builder.WriteByte(str[0] - ('a' - 'A'))
		index = 1
	}
	for i := index; i < len(str); i++ {
		if str[i] == '_' && i+1 < len(str) {
			if str[i+1] >= 'a' && str[i+1] <= 'z' {
				builder.WriteByte(str[i+1] - ('a' - 'A'))
				i++
				continue
			}
		}
		builder.WriteByte(str[i])
	}
	return builder.String()
}

/**
GetTimestamp 根据传入的时间字符串和格式获取时间戳，如果未传入则返回当前时间戳。
*/
func GetTimestamp(timeStr, layout string) (int64, error) {
	var parsedTime time.Time
	var err error

	if timeStr != "" && layout != "" {
		fmt.Println("按时间格式")
		// 解析时间字符串
		parsedTime, err = time.Parse(layout, timeStr)
		if err != nil {
			return 0, err
		}
	} else {
		// 没有传入时间字符串和格式，使用当前时间
		fmt.Println("wu时间格式")
		parsedTime = time.Now()
	}

	// 获取时间戳
	timestamp := parsedTime.Unix()
	return timestamp, nil
}

// IsTimeFormat 判断字符串是否是时间格式
func IsTimeFormat(str, layout string) bool {
	_, err := time.Parse(layout, str)
	return err == nil
}
