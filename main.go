// main.go
package main

import (
	"crypto/sha256"
	"encoding/base64"
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
	madeByLabel := widget.NewLabel("@xiaosheng : blog.devilwst.top")
	toolsLabel := widget.NewLabel("---小工具---")

	// 将自定义的 TextStyle 应用到标签的 TextStyle 属性上
	go views.StartClipboardListener(resultEntry)
	toolBtn := creatToolBtn(myApp)
	content := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("开启camel转换："),
		checkBox,
		resultEntry, // 添加文本框
		toolsLabel,
		toolBtn,
		madeByLabel,
	)

	listBox = widget.NewCheckGroup([]string{" 111"}, func(strings []string) {})
	listBox.Hide()
	sqlParseContent := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("sql输入："),
		container.NewHBox(widget.NewButton("输入 ", func() {
			views.SqlParsePre(myApp)
		}),
			widget.NewButton("解析", func() {
				ast.ParseSql(myApp)
				listBox.Show()
				flushColumnsToListBox(myWindow)
			}),
			widget.NewButton("生成", func() {
				tw := myApp.NewWindow("target")
				wet := widget.NewMultiLineEntry()
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
				flushColumnsToListBox(myWindow)
			}),
			listBox,
		))

	listBox = widget.NewCheckGroup(tools.SqlColumns, func(selected []string) {
		// 处理选择的选项
		fmt.Println("Selected:", selected)
		tools.SelectedRows = selected
	})
	csqlbox := container.New(layout.NewVBoxLayout(), sqlParseContent, listBox)
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

func GetTextWindow(myApp fyne.App) string {
	secret := ""
	sew := myApp.NewWindow("密钥")
	wne := widget.NewMultiLineEntry()
	wne.SetMinRowsVisible(8)
	wneBtn := widget.NewButton("确认", func() {
		secret = wne.Text
		sew.Close()
	})
	cb := container.NewVBox(wne, wneBtn)
	sew.SetContent(cb)
	sew.Show()
	return secret
}

func creatToolBtn(myApp fyne.App) fyne.CanvasObject {
	toolBtn := widget.NewButton("点击打开小工具", func() {
		tw := myApp.NewWindow("小工具")

		wetin := widget.NewMultiLineEntry()
		wetin.SetMinRowsVisible(10)
		wetout := widget.NewMultiLineEntry()
		wetout.SetMinRowsVisible(10)
		cbox := container.NewHBox(
			widget.NewButton("base64加密", func() {
				tools.Base64Origin = wetin.Text
				encoded := base64.StdEncoding.EncodeToString([]byte(tools.Base64Origin))
				fmt.Println("base64-encoded:", encoded)
				wetout.SetText(encoded)
			}),
			widget.NewButton("base64解密", func() {
				tools.Base64Encode = wetout.Text
				fmt.Println("base64-decoded-input:", wetout.Text)
				decoded, _ := base64.StdEncoding.DecodeString(tools.Base64Encode)
				fmt.Println("base64-decoded:", string(decoded))
				wetin.SetText(string(decoded))
			}),
			widget.NewButton("sha256加密", func() {
				fmt.Println("sha256-encoded-input:", wetin.Text)
				h := sha256.Sum256([]byte(wetin.Text))
				fmt.Println("base64-decoded:", string(h[:]))
				wetout.SetText(string(h[:]))
			}),
			widget.NewButton("md5", func() {
				tools.Base64Encode = wetout.Text
				fmt.Println("base64-decoded-input:", wetout.Text)
				decoded, _ := base64.StdEncoding.DecodeString(tools.Base64Encode)
				fmt.Println("base64-decoded:", string(decoded))
				wetin.SetText(string(decoded))
			}),
		)
		cn := container.NewVBox(wetin, cbox, wetout)
		tw.SetContent(cn)
		tw.Resize(fyne.NewSize(300, 300))
		tw.Show()

	})
	return toolBtn
}
