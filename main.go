// main.go
package main

import (
	"fmt"
	"xiaosheng/ast"
	"xiaosheng/tools"
	"xiaosheng/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

var listBox *widget.CheckGroup

func main() {
	myApp := app.New()
	icon, _ := fyne.LoadResourceFromPath("main.ico")
	myApp.SetIcon(icon)
	myWindow := myApp.NewWindow("剪贴板监听器")

	resultEntry := widget.NewEntry()
	resultEntry.MultiLine = true
	resultEntry.Disable()
	// 创建复选框
	checkBox := widget.NewCheck("监听剪贴板", func(value bool) {
		fmt.Println("flag:", value)
		views.ShouldListenClipboard = value // 设置标志来表示是否要监听剪贴板
	})

	// 创建一个标签
	madeByLabel := widget.NewLabel("@xiaosheng ")
	toolsLabel := widget.NewLabel("---小工具---")

	// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
	go views.StartClipboardListener(resultEntry)
	toolBtn := views.CreatToolBtn(myApp)
	content := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("开启camel转换："),
		checkBox,
		resultEntry, // 添加文本框
		toolsLabel,
		toolBtn,
		madeByLabel,
	)

	sqlParseContent := SqlContent(myApp, &myWindow)

	csqlbox := container.New(layout.NewVBoxLayout(), sqlParseContent)
	csqlbox.Resize(fyne.NewSize(300, 300))

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), content, csqlbox))
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.Content().Size().Max(fyne.NewSize(1920, 1000))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}

func flushColumnsToListBox(myWindow fyne.Window) {
	listBox.Options = tools.SqlColumns
	myWindow.Content().Refresh()
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
			views.SqlParsePre(myApp)
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
