// views/listen_clipboard.go
package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/flopp/go-findfont"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"xiaosheng/logs"
	"xiaosheng/tools"
)

var (
	ShouldListenClipboard bool // 新增标志来表示是否要监听剪贴板
	ShouldLog             bool // 新增标志来表示是否要监听剪贴板
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

func StartClipboardListener(resultEntry *widget.Entry, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			if ShouldListenClipboard {
				// 读取剪贴板内容
				clipboardText, _ := clipboard.ReadAll()
				originText := clipboardText

				// 判断剪贴板内容是否是英文或下划线连接
				if isEnglishOrUnderscore(clipboardText) {
					// 如果是下划线连接的字符串，将其转换为驼峰命名
					camelCase := tools.ToCamelCase(clipboardText)

					if originText != camelCase {
						// 将转换后的内容复制回剪贴板
						clipboard.WriteAll(camelCase)
						// 更新文本框的内容
						resultEntry.SetText("监听到：" + clipboardText + "\n替换为：" + camelCase)
					}

				} else {
					resultEntry.SetText("工作啦伙计") // 清空文本框
				}
			} else {
				//time.Sleep(300 * time.Millisecond) // 不监听剪贴板时，降低CPU负载
				resultEntry.SetText("休眠中。。。") // 清空文本框
			}
		}
	}
}

func ListenClipBordPart(myApp fyne.App) *fyne.Container {
	resultEntry := widget.NewEntry()
	resultEntry.MultiLine = true
	resultEntry.Disable()
	// 创建复选框
	ticker := time.NewTicker(200 * time.Millisecond)
	logTicker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	checkBox := widget.NewCheck("监听剪贴板", func(value bool) {
		log.Println("监听 剪切板：flag:", value)
		ShouldListenClipboard = value // 设置标志来表示是否要监听剪贴板
		if value {
			// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
			ticker = time.NewTicker(70 * time.Millisecond)
			go StartClipboardListener(resultEntry, ticker)
		} else {
			log.Println("关闭监听剪贴板")
			ticker.Stop()
		}
	})
	logCheckBox := widget.NewCheck("isLog", func(value bool) {
		log.Println("log：flag:", value)
		ShouldLog = value // 设置标志来表示是否要监听剪贴板
		if value {
			// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
			// 初始化日志
			logs.SetupLogger()
			logTicker = time.NewTicker(5 * time.Second)
			go logs.MonitorFileSize(200*1024*1024, logTicker)
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

	toolBtn := CreatToolBtn(myApp)
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
