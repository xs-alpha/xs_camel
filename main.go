package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/flopp/go-findfont"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	shouldListenClipboard bool // 新增标志来表示是否要监听剪贴板
)

func init() {
	// 设置中文字体
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func isEnglishOrUnderscore(text string) bool {
	// 判断字符串是否是英文或下划线连接
	matched, _ := regexp.MatchString("^[A-Za-z_]+$", text)
	return matched
}

func toCamelCase(text string) string {
	// 将下划线连接的字符串转换为驼峰命名
	parts := strings.Split(text, "_")
	camelCase := parts[0]
	for i := 1; i < len(parts); i++ {
		camelCase += strings.Title(parts[i])
	}
	return camelCase
}

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
		shouldListenClipboard = value // 设置标志来表示是否要监听剪贴板
	})

	go func() {
		for {
			if shouldListenClipboard {
				// 读取剪贴板内容
				clipboardText, _ := clipboard.ReadAll()
				originText := clipboardText

				// 判断剪贴板内容是否是英文或下划线连接
				if isEnglishOrUnderscore(clipboardText) {
					// 如果是下划线连接的字符串，将其转换为驼峰命名
					camelCase := toCamelCase(clipboardText)

					if originText != camelCase {
						// 将转换后的内容复制回剪贴板
						clipboard.WriteAll(camelCase)
						// 更新文本框的内容
						resultEntry.SetText("监听到：" + clipboardText + "\n替换为：" + camelCase)
					}

				} else {
					resultEntry.SetText("") // 清空文本框
				}

				// 等待一段时间再继续检查剪贴板
				time.Sleep(200 * time.Millisecond)
			} else {
				time.Sleep(300 * time.Millisecond) // 不监听剪贴板时，降低CPU负载
				resultEntry.SetText("休眠中。。。")      // 清空文本框
			}
		}
	}()
	// 创建一个标签
	madeByLabel := widget.NewLabel("@xiaosheng : blog.devilwst.top")

	// 创建一个自定义的 TextStyle 结构，并设置字体大小
	customTextStyle := fyne.TextStyle{
		Bold:      false,
		Italic:    true,
		Monospace: false, // 可根据需要设置其他样式
	}
	// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
	madeByLabel.TextStyle = customTextStyle

	content := container.NewVBox(
		widget.NewLabel("选择是否监听剪贴板："),
		checkBox,
		resultEntry, // 添加文本框
		madeByLabel,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}
