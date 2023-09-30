package main

import (
"fmt"
"io/ioutil"
"regexp"
"strings"

	"github.com/xwb1989/sqlparser"
)

var sqlTypeMap = map[string]string{
"int":                "Integer",
"integer":            "Integer",
"tinyint":            "Integer",
"smallint":           "Integer",
"mediumint":          "Integer",
"bigint":             "Integer",
"int unsigned":       "Integer",
"integer unsigned":   "Integer",
"tinyint unsigned":   "Integer",
"smallint unsigned":  "Integer",
"mediumint unsigned": "Integer",
"bigint unsigned":    "Integer",
"bit":                "Byte",
"bool":               "Boolean",
"enum":               "String",
"set":                "String",
"varchar":            "String",
"char":               "String",
"tinytext":           "String",
"mediumtext":         "String",
"text":               "String",
"longtext":           "String",
"blob":               "String",
"tinyblob":           "String",
"mediumblob":         "String",
"longblob":           "String",
"date":               "date",
"datetime":           "datetime",
"timestamp":          "timestamp",
"time":               "time",
"float":              "float64",
"double":             "float64",
"decimal":            "decimal",
"binary":             "binary",
"varbinary":          "binary",
}

// SqlToGo converts a sql create statement to Go struct
// sqlStmt for sql create statement, outputPkg for output directory
func SqlToGo(sqlStmt string, outputPkg string) (string, error) {
statement, err := sqlparser.ParseStrictDDL(sqlStmt)
if err != nil {
fmt.Println("errFunc")
return "", err
}
stmt, ok := statement.(*sqlparser.DDL)
if !ok {
return "", fmt.Errorf("input sql is not a create statment")
}
// convert to Go struct
tableName := stmt.NewName.Name.String()
fmt.Println("convert to go---------------")
res, err := stmtToGo(stmt, tableName, outputPkg)
if err != nil {
return "", err
}
return res, nil
}

func stmtToGo(stmt *sqlparser.DDL, tableName string, pkgName string) (string, error) {
builder := strings.Builder{}

	header := fmt.Sprintf("package %s\n", pkgName)

	structName := snakeCaseToCamel(tableName)
	structStart := fmt.Sprintf("type %s struct { \n", structName)
	builder.WriteString(structStart)
	for _, col := range stmt.TableSpec.Columns {
		columnType := col.Type.Type


		goType := sqlTypeMap[columnType]

		field := snakeCaseToCamel(col.Name.String())
		comment := col.Type.Comment
		if comment == nil {
			builder.WriteString(fmt.Sprintf("\t%s\t%s\t\n", field, goType))
		} else {
			builder.WriteString(fmt.Sprintf("\t%s\t%s\t`comment:\"%s\"` \n",
				field, goType, string(comment.Val)))
		}
	}
	builder.WriteString("}\n")

	return header + builder.String(), nil
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

func main() {
data, err := ioutil.ReadFile("./input.sql")
if err != nil {
fmt.Println("读取err:",err)
return
}
sqlStatement := string(data)
// Remove parentheses after 'timestamp' and 'CURRENT_TIMESTAMP' (case-insensitive)
re := regexp.MustCompile(`(?i)(timestamp|current_timestamp)\(\d+\)`)
sqlStatement = re.ReplaceAllString(sqlStatement, "$1")

	// Remove 'ON UPDATE CURRENT_TIMESTAMP' and the following parentheses (case-insensitive)
	re = regexp.MustCompile(`(?i)ON UPDATE CURRENT_TIMESTAMP\(\d+\)`)
	sqlStatement = re.ReplaceAllString(sqlStatement, "")

	
	fmt.Println("sql:  ->",sqlStatement)

	res,err := SqlToGo(sqlStatement, "xiaosheng")
	if err != nil {
		fmt.Println("解析err:"+err.Error())
		return
	}
	// print result
	fmt.Println("res"+res)
}


 err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "qr.png")