// main.go
package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"xiaosheng/ast"
	"xiaosheng/tools"
	"xiaosheng/views"

	"fyne.io/fyne/v2/canvas"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	qrcodeGenerate "github.com/skip2/go-qrcode"
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
	listBoxContainer := container.NewVScroll(listBox) // Wrap the CheckGroup in a scrollable container
	listBoxContainer.Hide()
	sqlParseContent := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("sql输入："),
		container.NewHBox(widget.NewButton("输入 ", func() {
			views.SqlParsePre(myApp)
		}),
			widget.NewButton("解析", func() {
				ast.ParseSql(myApp)
				fmt.Println("tools.column:", tools.SqlColumns)
				listBox.Show()
				listBoxContainer.Show()
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
		listBoxContainer,
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

func GetTextWindow(myApp fyne.App, title string, label string) string {
	secret := ""
	sew := myApp.NewWindow(title)
	wn := widget.NewLabel(label)
	wne := widget.NewMultiLineEntry()
	wne.SetMinRowsVisible(8)
	wneBtn := widget.NewButton("确认", func() {
		secret = wne.Text
		sew.Close()
	})
	cb := container.NewVBox(wn,wne, wneBtn)
	sew.SetContent(cb)
	sew.Resize(fyne.NewSize(200, 200))
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
				dia := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
					if err != nil {
						fmt.Println("newOpenFileFolder:", err.Error())
						return
					}
					suffix := uc.URI().Extension()
					fmt.Println("suffix:", suffix)
					if !tools.IsImg(suffix) {
						wetin.SetText("请打开图片格式文件哦！")
						fmt.Println("notImg-suffix:", suffix)
						return
					}
					path := uc.URI().String()
					wetin.SetText("打开二维码路径：" + path)
					fmt.Println("path:", path)
					path = strings.Replace(path, "file://", "", 1)
					res := tools.ReadQRCode(path)
					fmt.Println("res:", res)
					wetout.SetText("解析 成功-- 解析内容为： " + res)
				}, tw)
				dia.Show()
			}),

			widget.NewButton("二维码生成", func() {
				//cont:=GetTextWindow(myApp,"输入二维码内容","请在输入框输入二维码 内容：")
				sew := myApp.NewWindow("输入二维码内容")
				wn := widget.NewLabel("请在输入框输入二维码 内容：")
				wne := widget.NewMultiLineEntry()
				wne.SetMinRowsVisible(8)
				wneBtn := widget.NewButton("确认", func() {
					cont := wne.Text
					sew.Close()
					var png []byte
					png, err := qrcodeGenerate.Encode(cont, qrcodeGenerate.Medium, 256)
					if err != nil {
						wetin.SetText("生成 二维码失败")
					}
					pngReader := bytes.NewReader(png)
					qrwin := myApp.NewWindow("生成二维码")
					image := canvas.NewImageFromReader(pngReader, "QR")
					image.FillMode = canvas.ImageFillOriginal
					qrctn := container.NewVBox(image, widget.NewButton("保存", func() {
						if isSavedSuccess, filename := writeInFile(pngReader); isSavedSuccess {
							wetout.SetText("二维码保存成功,文件名：" + filename)
							qrwin.Close()
						} else {
							wetout.SetText("二维码 保存失败")
						}
					}))
					qrwin.SetContent(qrctn)
					qrwin.Resize(fyne.NewSize(200, 200))
					qrwin.Show()
				})
				cb := container.NewVBox(wn,wne, wneBtn)
				sew.SetContent(cb)
				sew.Resize(fyne.NewSize(200, 200))
				sew.Show()
				//err := qrcodeGenerate.WriteFile("https://example.org", qrcodeGenerate.Medium, 256, "qr.png")
			}),
			widget.NewButton("sha512加密", func() {
				fmt.Println("sha512-input:", wetin.Text)
				hashedBytes := sha512.Sum512([]byte(wetin.Text))
				hashedString := fmt.Sprintf("%x", hashedBytes)
				fmt.Println("sha512-output:", hashedString)
				wetout.SetText(hashedString)
			}),
			widget.NewButton("随机密码", func() {
				wetin.PlaceHolder = "可以输入要生成的随机密码长度，不输入默认取8"
				pwd:=""
				if !tools.IsNumeric(wetin.Text){
					pwd,_ = tools.GenerateRandomPassword(8)
				}else{
					leng,_:=strconv.Atoi(wetin.Text)
					pwd,_=tools.GenerateRandomPassword(leng)
				}
				wetout.SetText("生成随记密码:"+ pwd)
				fmt.Println("password-output:", pwd)
			}),
			widget.NewButton("json美化", func() {
				// 解析JSON数据到一个map或结构体
				var parsedData map[string]interface{}
				if err := json.Unmarshal([]byte(wetin.Text), &parsedData); err != nil {
					fmt.Println("JSON解析失败:", err)
					wetout.SetText("json解析失败："+err.Error())
					return
				}
				prettyJson,err:=tools.PrettyPrintJSON(parsedData)
				if err!=nil{
					wetout.SetText("json解析失败")
					return
				}
				fmt.Println("json美化："+prettyJson)
				wetout.SetText(prettyJson)
			}),
			widget.NewButton("字数统计", func() {
				fmt.Println("字数统计-input:", wetin.Text)
				num:=tools.CountValidWords(wetin.Text)
				fmt.Println("字数统计-output:", num)
				wetout.SetText("字数统计(不统计空格):"+strconv.Itoa(num))
			}),
		)
		cn := container.NewVBox(wetin, cbox, cboxImg, wetout)
		tw.SetContent(cn)
		tw.Resize(fyne.NewSize(300, 300))
		tw.Show()

	})
	return toolBtn
}

func readImgIO(f fyne.URIReadCloser) {
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

func writeInFile(imageReader io.Reader) (bool, string) {
	// 获取当前时间戳并格式化为字符串
	timestamp := time.Now().Unix()
	timestampStr := fmt.Sprintf("%d", timestamp)

	// 创建文件名，将时间戳作为文件名的一部分，加上文件后缀（例如 ".png"）
	fileName := timestampStr + ".png"
	// 创建一个文件用于保存图片
	outputFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return false, ""
	}
	defer outputFile.Close()

	// 将 imageReader 重新定位到开头
	if _, err := imageReader.(*bytes.Reader).Seek(0, 0); err != nil {
		fmt.Println("Error seeking imageReader:", err)
		return false, ""
	}

	// 使用 io.Copy 将数据从 io.Reader 复制到文件
	_, copyErr := io.Copy(outputFile, imageReader)
	if copyErr != nil {
		fmt.Println("Error copying data to file:", copyErr)
		return false, ""
	}

	fmt.Println("Image saved as " + fileName)
	return true, fileName
}
