// main.go
package main

import (
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
	"xiaosheng/logs"
	"xiaosheng/tools"
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
	content := ListenClipBordPart(&myApp)
	// 右侧sql解析
	sqlParseContent := views.SqlContent(myApp, &myWindow)

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), content, sqlParseContent))
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.Content().Size().Max(fyne.NewSize(1920, 1000))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}

func ListenClipBordPart(myApp *fyne.App) *fyne.Container {
	resultEntry := widget.NewEntry()
	resultEntry.MultiLine = true
	resultEntry.Disable()
	// 创建复选框
	ticker := time.NewTicker(200 * time.Millisecond)
	logTicker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	checkBox := widget.NewCheck("监听剪贴板", func(value bool) {
		log.Println("监听 剪切板：flag:", value)
		views.ShouldListenClipboard = value // 设置标志来表示是否要监听剪贴板
		if value {
			// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
			ticker = time.NewTicker(70 * time.Millisecond)
			go views.StartClipboardListener(resultEntry, ticker)
			defer ticker.Stop()
		} else {
			log.Println("关闭监听剪贴板")
			ticker.Stop()
		}
	})
	logCheckBox := widget.NewCheck("isLog", func(value bool) {
		log.Println("log：flag:", value)
		views.ShouldLog = value // 设置标志来表示是否要监听剪贴板
		if value {
			// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
			// 初始化日志
			logs.SetupLogger()
			logTicker = time.NewTicker(5 * time.Second)
			go logs.MonitorFileSize(200*1024*1024, logTicker)
			defer ticker.Stop()
		} else {
			log.Println("关闭日志")
			logs.CloseLogger()
			logTicker.Stop()
		}
	})
	camelBox := widget.NewCheck("大驼峰", func(value bool) {
		log.Println("大驼峰：flag:", value)
		tools.IsBigCamel = value // 设置标志来表示是否要监听剪贴板
	})

	checkBoxContainer := container.NewHBox(checkBox, camelBox)

	// 创建一个标签
	madeByLabel := widget.NewLabel("	  @xiaosheng 	 ")
	toolsLabel := widget.NewLabel("	   小工具")

	toolBtn := views.CreatToolBtn(*myApp)
	content := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("开启camel转换："),
		checkBoxContainer,
		resultEntry, // 添加文本框
		toolsLabel,
		toolBtn,
		madeByLabel,
		logCheckBox,
	)
	return content
}
