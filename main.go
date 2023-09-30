// main.go
package main

import (
	"fmt"
	"xiaosheng/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	icon, _ := fyne.LoadResourceFromPath("main.ico")
	myApp.SetIcon(icon)
	myWindow := myApp.NewWindow("小生 开发助手v0.3")

	resultEntry := widget.NewEntry()
	resultEntry.MultiLine = true
	resultEntry.Disable()
	// 创建复选框
	checkBox := widget.NewCheck("监听剪贴板", func(value bool) {
		fmt.Println("flag:", value)
		views.ShouldListenClipboard = value // 设置标志来表示是否要监听剪贴板
	})

	// 创建一个标签
	madeByLabel := widget.NewLabel("	  @xiaosheng 	 ")
	toolsLabel := widget.NewLabel("	   小工具")

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
	content.Resize(fyne.NewSize(200, 300))

	sqlParseContent := views.SqlContent(myApp, &myWindow)

	csqlbox := container.New(layout.NewVBoxLayout(), sqlParseContent)
	csqlbox.Resize(fyne.NewSize(300, 300))

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), content, csqlbox))
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.Content().Size().Max(fyne.NewSize(1920, 1000))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}
