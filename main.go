// main.go
package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"xiaosheng/tools"
	"xiaosheng/views"
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
	madeByLabel := widget.NewLabel("@xiaosheng : blog.devilwst.top")

	// 创建一个自定义的 TextStyle 结构，并设置字体大小
	customTextStyle := fyne.TextStyle{
		Bold:      false,
		Italic:    true,
		Monospace: false, // 可根据需要设置其他样式
	}
	// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
	go views.StartClipboardListener(resultEntry)
	madeByLabel.TextStyle = customTextStyle

	content := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("选择是否监听剪贴板："),
		checkBox,
		resultEntry, // 添加文本框
		madeByLabel,
	)
	content.Resize(fyne.NewSize(200, 200))

	listBox = widget.NewCheckGroup([]string{" 111"}, func(strings []string) {})
	listBox.Hide()
	sqlParseContent := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("sql输入："),
		container.NewHBox(widget.NewButton("输入 ", func() {
			views.SqlParsePre(myApp)
		}),
			widget.NewButton("解析", func() {
				views.ParseSql(myApp)
				listBox.Show()
				flushColumnsToListBox(myWindow)
			}),
			widget.NewButton("生成", func() {
				tw := myApp.NewWindow("target")
				wet := widget.NewMultiLineEntry()
				element := tools.GetExcelElement(views.SelectedRows)
				wet.SetText(element)
				wet.SetMinRowsVisible(20)

				box := container.NewVBox(wet, widget.NewButton("复制", func() {

				}))
				tw.SetContent(box)
				tw.Resize(fyne.NewSize(300, 300))
				tw.Show()
			}),
			widget.NewButton("clear", func() {
				views.SqlColumns = []string{}
				views.SqlStatement = ""
				flushColumnsToListBox(myWindow)
			}),
			listBox,
		))

	listBox = widget.NewCheckGroup(views.SqlColumns, func(selected []string) {
		// 处理选择的选项
		fmt.Println("Selected:", selected)
		views.SelectedRows = selected
	})
	csqlbox := container.New(layout.NewVBoxLayout(), sqlParseContent, listBox)
	csqlbox.Resize(fyne.NewSize(300, 300))

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), content, csqlbox))
	myWindow.Resize(fyne.NewSize(500, 300))

	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}

func flushColumnsToListBox(myWindow fyne.Window) {
	listBox.Options = views.SqlColumns
	myWindow.Content().Refresh()
}
