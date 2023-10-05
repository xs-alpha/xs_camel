// views/listen_clipboard.go
package views

import (
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/flopp/go-findfont"
	"os"
	"regexp"
	"strings"
	"time"
	"xiaosheng/tools"
)

var (
	ShouldListenClipboard bool // 新增标志来表示是否要监听剪贴板
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

func StartClipboardListener(resultEntry *widget.Entry) {
	go func() {
		for {
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

				// 等待一段时间再继续检查剪贴板
				time.Sleep(200 * time.Millisecond)
			} else {
				time.Sleep(300 * time.Millisecond) // 不监听剪贴板时，降低CPU负载
				resultEntry.SetText("休眠中。。。")      // 清空文本框
			}
		}
	}()
}
