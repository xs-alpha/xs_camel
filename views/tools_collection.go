package views

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"xiaosheng/tools"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	qrcodeGenerate "github.com/skip2/go-qrcode"
)

func CreatToolBtn(myApp fyne.App) fyne.CanvasObject {
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
						if isSavedSuccess, filename := tools.WriteInFile(pngReader); isSavedSuccess {
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
				cb := container.NewVBox(wn, wne, wneBtn)
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
				pwd := ""
				if !tools.IsNumeric(wetin.Text) {
					pwd, _ = tools.GenerateRandomPassword(8)
				} else {
					leng, _ := strconv.Atoi(wetin.Text)
					pwd, _ = tools.GenerateRandomPassword(leng)
				}
				wetout.SetText("生成随记密码:" + pwd)
				fmt.Println("password-output:", pwd)
			}),
			widget.NewButton("json美化", func() {
				// 解析JSON数据到一个map或结构体
				var parsedData map[string]interface{}
				if err := json.Unmarshal([]byte(wetin.Text), &parsedData); err != nil {
					fmt.Println("JSON解析失败:", err)
					wetout.SetText("json解析失败：" + err.Error())
					return
				}
				prettyJson, err := tools.PrettyPrintJSON(parsedData)
				if err != nil {
					wetout.SetText("json解析失败")
					return
				}
				fmt.Println("json美化：" + prettyJson)
				wetout.SetText(prettyJson)
			}),
			widget.NewButton("字数统计", func() {
				fmt.Println("字数统计-input:", wetin.Text)
				num := tools.CountValidWords(wetin.Text)
				fmt.Println("字数统计-output:", num)
				wetout.SetText("字数统计(不统计空格):" + strconv.Itoa(num))
			}),
		)
		cn := container.NewVBox(wetin, cbox, cboxImg, wetout)
		tw.SetContent(cn)
		tw.Resize(fyne.NewSize(300, 300))
		tw.Show()

	})
	return toolBtn
}
