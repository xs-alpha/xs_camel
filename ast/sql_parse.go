package ast

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/xwb1989/sqlparser"
	"regexp"
	"xiaosheng/tools"
)

// SqlToGo converts a sql create statement to Go struct
// sqlStmt for sql create statement, outputPkg for output directory
func SqlToGo(sqlStmt string, outputPkg string) ([]string, error) {
	statement, err := sqlparser.ParseStrictDDL(sqlStmt)
	if err != nil {
		fmt.Println("errFunc")
		return []string{}, err
	}
	stmt, ok := statement.(*sqlparser.DDL)
	if !ok {
		return []string{}, fmt.Errorf("input sql is not a create statment")
	}
	// convert to Go struct
	tableName := stmt.NewName.Name.String()
	fmt.Println("convert to go---------------")
	res, err := tools.StmtToGo(stmt, tableName, outputPkg)
	if err != nil {
		return []string{}, err
	}
	return res, nil
}

func ParseSql(app fyne.App) {
	//data, err := ioutil.ReadFile("./input.sql")
	//if err != nil {
	//	fmt.Println("读取err:",err)
	//	return
	//}
	if tools.SqlStatement == "" {
		window := app.NewWindow("info")
		label := widget.NewLabel("请输入正确的sql语句")
		window.SetContent(container.New(layout.NewVBoxLayout(), label, widget.NewButton("确认", func() {
			window.Close()
		})))
		window.Show()

	}
	sqlStatement := string(tools.SqlStatement)
	// Remove parentheses after 'timestamp' and 'CURRENT_TIMESTAMP' (case-insensitive)
	re := regexp.MustCompile(`(?i)(timestamp|current_timestamp)\(\d+\)`)
	sqlStatement = re.ReplaceAllString(sqlStatement, "$1")

	// Remove 'ON UPDATE CURRENT_TIMESTAMP' and the following parentheses (case-insensitive)
	re = regexp.MustCompile(`(?i)ON UPDATE CURRENT_TIMESTAMP\(\d+\)`)
	sqlStatement = re.ReplaceAllString(sqlStatement, "")

	fmt.Println("sql:  ->", sqlStatement)
	if sqlStatement == "" {
		return
	}

	res, err := SqlToGo(sqlStatement, "xiaosheng")
	if err != nil {
		fmt.Println("解析err:" + err.Error())
		return
	}
	// print result
	for _, v := range res {
		fmt.Println("res" + v)
	}
	// 遍历原始切片，筛选出非空字符串并添加到新的切片中
	//tools.SqlColumns = res

	if tools.IsAppended {
		l := len(tools.SqlColumns)
		// 遍历原始切片并为每个元素前面加上索引加上点号
		for i, value := range res {
			res[i] = fmt.Sprintf("%d.%s", i+l, value)
		}

		tools.SqlColumns = append(tools.SqlColumns, res...)
		fmt.Println("all:", tools.SqlColumns, "now:", res)
	} else {
		tools.SqlColumns = res
		// 遍历原始切片并为每个元素前面加上索引加上点号
		for i, value := range tools.SqlColumns {
			tools.SqlColumns[i] = fmt.Sprintf("%d.%s", i, value)
		}

	}

}
