package views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"xiaosheng/ast"
	"xiaosheng/tools"
)

var listBox *widget.CheckGroup

func SqlParsePre(myApp fyne.App) {
	w3 := myApp.NewWindow("sql字串输入")
	we := widget.NewMultiLineEntry()
	we.SetMinRowsVisible(15) // 设置高度为 400
	sqlBox := container.NewVBox(we)
	//sqlBox.Resize(fyne.NewSize(300, 400)) // 设置高度为 400
	button := widget.NewButton("确认", func() {
		tools.SqlStatement = we.Text
		fmt.Println("监听到sql", we.Text)
		w3.Close()
	})

	wec := container.New(layout.NewVBoxLayout(), widget.NewLabel("请输入sql 建表语句"), sqlBox, button)
	w3.SetContent(wec)
	w3.Resize(fyne.NewSize(300, 400))
	w3.Show()
}

// sql解析部分ui
func SqlContent(myApp fyne.App, myWindow *fyne.Window) *fyne.Container {
	// myWindow:=*myWindo
	listBox = widget.NewCheckGroup([]string{" 111"}, func(selected []string) {
		fmt.Println("Selected:", selected)
		tools.SelectedRows = selected
	})
	//listBox.Hide()
	listBoxContainer := container.NewVScroll(listBox) // Wrap the CheckGroup in a scrollable container
	listBoxContainer.Resize(fyne.NewSize(200, 400))
	listBoxContainer.SetMinSize(fyne.NewSize(200, 350))
	listBoxContainer.Hide()
	sqlParseContent := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("sql输入："),
		container.NewHBox(widget.NewButton("输入 ", func() {
			SqlParsePre(myApp)
		}),
			widget.NewButton("解析", func() {
				ast.ParseSql(myApp)
				fmt.Println("tools.column:", tools.SqlColumns)
				flushColumnsToListBox(*myWindow)
				//listBox.Show()
				listBoxContainer.Show()
				listBoxContainer.Refresh()
				listBox.Refresh()
			}),
			widget.NewButton("生成", func() {
				tw := myApp.NewWindow("target")
				wet := widget.NewMultiLineEntry()
				fmt.Println("tools.selected:", tools.SelectedRows)
				element := tools.GetExcelElement(tools.SelectedRows)
				wet.SetText(element)
				wet.SetMinRowsVisible(20)
				tools.SelectedText = element
				box := container.NewVBox(wet, widget.NewButton("复制", func() {
					clipboard.WriteAll(tools.SelectedText)
					tw.Close()
				}))
				tw.SetContent(box)
				tw.Resize(fyne.NewSize(300, 300))
				tw.Show()
			}),
			widget.NewButton("clear", func() {
				tools.SqlColumns = []string{}
				tools.SqlStatement = ""
				tools.SelectedRows = []string{}
				// 清空 listBox 中的选定项
				listBox.SetSelected([]string{})
				flushColumnsToListBox(*myWindow)
				listBox.Refresh()
				listBoxContainer.Refresh()

			}),
			widget.NewCheck("是否追加", func(val bool) {
				tools.IsAppended = val
			}),
		),
		//listBoxContainer,
		//listBox,
	)
	tmp := container.NewVBox(sqlParseContent, listBoxContainer)
	tmp.Resize(fyne.NewSize(200, 400))
	return tmp
	//return sqlParseContent
}
func flushColumnsToListBox(myWindow fyne.Window) {
	listBox.Options = tools.SqlColumns
	myWindow.Content().Refresh()
}
