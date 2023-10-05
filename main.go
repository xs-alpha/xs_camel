// main.go
package main

import (
	"xiaosheng/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func main() {
	myApp := app.New()
	icon, _ := fyne.LoadResourceFromPath("main.ico")
	myApp.SetIcon(icon)
	myWindow := myApp.NewWindow("小生 开发助手v0.6")

	// 左侧监听剪贴板
	content := views.ListenClipBordPart(myApp)
	// 右侧sql解析
	sqlParseContent := views.SqlContent(myApp, &myWindow)

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), content, sqlParseContent))
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.Content().Size().Max(fyne.NewSize(1920, 1000))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}
