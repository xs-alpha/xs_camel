// main.go
package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"xiaosheng/ast"
	"xiaosheng/tools"
	"xiaosheng/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	listBox = widget.NewCheckGroup([]string{" 111"}, func(selected []string) {
		fmt.Println("Selected:", selected)
		tools.SelectedRows = selected
	})
	listBox.Hide()
	sqlParseContent := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("sql输入："),
		container.NewHBox(widget.NewButton("输入 ", func() {
			views.SqlParsePre(myApp)
		}),
			widget.NewButton("解析", func() {
				ast.ParseSql(myApp)
				fmt.Println("tools.column:", tools.SqlColumns)
				listBox.Show()
				flushColumnsToListBox(myWindow)
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
				flushColumnsToListBox(myWindow)
			}),
			widget.NewCheck("是否追加", func(val bool) {
				tools.IsAppended = val
			}),
		),
		listBox,
	)

	//listBox = widget.NewCheckGroup(tools.SqlColumns, func(selected []string) {
	//	// 处理选择的选项
	//	fmt.Println("Selected:", selected)
	//	tools.SelectedRows = selected
	//})

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
		wetin.PlaceHolder = "原文"
		wetout := widget.NewMultiLineEntry()
		wetout.SetMinRowsVisible(10)
		wetout.PlaceHolder = "密文"
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
				fmt.Println("base64-decoded:", fmt.Sprintf("%x", h))
				wetout.SetText(fmt.Sprintf("%x", h))
			}),
			widget.NewButton("md5", func() {
				fmt.Println("md5-input:", wetin.Text)
				sum := md5.Sum([]byte(wetin.Text))
				fmt.Println("md5-output:", fmt.Sprintf("%x", sum))
				wetout.SetText(fmt.Sprintf("%x", sum))
			}),
			widget.NewButton("urlEncode", func() {
				fmt.Println("urlEncode-input:", wetin.Text)
				sum := url.QueryEscape(wetin.Text)
				fmt.Println("urlEncode-output:", sum)
				wetout.SetText(sum)
			}),
			widget.NewButton("urlDecode", func() {
				fmt.Println("urlDecode-input:", wetout.Text)
				out, _ := url.QueryUnescape(wetout.Text)
				fmt.Println("urlDecode-output:", out)
				wetout.SetText(out)
			}),
			widget.NewButton("时间戳", func() {
				isTime := tools.IsTimeFormat(wetin.Text, "2006-01-02 15:04:05")
				if !isTime {
					wetin.SetText("")
				}
				wetin.PlaceHolder = "传入时间格式yyyy-MM-dd HH:mm:ss可输出某时刻的时间戳，未传入这输出当前时间戳!!"
				fmt.Println("时间戳-input:", wetin.Text)
				timestamp, _ := tools.GetTimestamp(wetin.Text, "2006-01-02 15:04:05")
				fmt.Println("时间戳-output:", strconv.FormatInt(timestamp, 10))
				if isTime {
					wetout.SetText("时间戳： " + strconv.FormatInt(timestamp, 10))
				} else {
					wetout.SetText("当前时间戳： " + strconv.FormatInt(timestamp, 10))
				}
			}),
		)
		cboxImg := container.NewHBox(
			widget.NewButton("二维码解码", func() {
				// wetin.Hide()

				dia:=dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
					if err!=nil{
						fmt.Println("newOpenFileFolder:",err.Error())
						return
					}
					suffix:=uc.URI().Extension()
					fmt.Println("suffix:",suffix)
					if !tools.IsImg(suffix){
						fmt.Println("notImg-suffix:",suffix)
						return
					}
					path:=uc.URI().String()
					wetin.SetText("打开二维码路径："+path)
					fmt.Println("path:",path)
					tools.ReadQRCode(path)
				}, tw)
				dia.Show()
				fmt.Println("urlDecode-input:", wetout.Text)
				out, _ := url.QueryUnescape(wetout.Text)
				fmt.Println("urlDecode-output:", out)
				wetout.SetText(out)
			}),
		)
		cn := container.NewVBox(wetin, cbox,cboxImg, wetout)
		tw.SetContent(cn)
		tw.Resize(fyne.NewSize(300, 300))
		tw.Show()

	})
	return toolBtn
}

func readImgIO(f fyne.URIReadCloser){
	full := make([]byte, 0)
	b := make([]byte, 1024)
	for {
		_, err := f.Read(b)
		full = append(full, b...)
		if err == io.EOF {
			break
		}
	}
	fmt.Println(full)
	defer f.Close()
}