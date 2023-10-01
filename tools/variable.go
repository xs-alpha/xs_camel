package tools

import "fyne.io/fyne/v2/widget"

var (
	SqlStatement string
	SqlColumns   []string
	SelectedRows []string
	SelectedText string
	IsAppended   bool
	ListBox      *widget.CheckGroup
	IsBigCamel   bool

	// base64
	Base64Origin string
	Base64Encode string
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
