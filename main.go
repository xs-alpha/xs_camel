// main.go
package main

import (
	"flag"
	"fmt"
	"xiaosheng/settings"
	"xiaosheng/views"
	"xiaosheng/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func main() {
	// 初始化配置文件
	var file string
	var size int
	flag.IntVar(&size, "size", 200, "文件大小")
	flag.StringVar(&file, "configFile", "config.yaml", "配置文件")
	flag.Parse()
	if err := settings.Init(file); err != nil {
		fmt.Printf("main: viper initiallize failed :%v\n", err.Error())
		return
	}
	// 配置窗口
	myApp := app.New()
	myApp.Settings().SetTheme(&theme.MyTheme{})
	icon, _ := fyne.LoadResourceFromPath("main.ico")
	myApp.SetIcon(icon)
	myWindow := myApp.NewWindow("小生 开发助手v" + settings.Conf.SoftVersion)

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
