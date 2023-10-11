package views

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"xiaosheng/settings"
	"xiaosheng/tools"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	qrcodeGenerate "github.com/skip2/go-qrcode"
)

var wg  sync.WaitGroup

func CreatToolBtn(myApp fyne.App) fyne.CanvasObject {
	toolBtn := widget.NewButton("点击打开小工具", func() {
		tw := myApp.NewWindow("小生开发助手——小工具v" + settings.Conf.SoftVersion)

		wetin := widget.NewMultiLineEntry()
		wetin.SetMinRowsVisible(10)
		wetin.PlaceHolder = "原文"
		wetout := widget.NewMultiLineEntry()
		wetout.SetMinRowsVisible(10)
		wetout.PlaceHolder = "result"
		wetres := widget.NewMultiLineEntry()
		wetres.SetMinRowsVisible(4)
		wetres.PlaceHolder = "result"

		ticker := time.NewTicker(100 * time.Millisecond)
		tools.ToolsChan = make(chan int, 1)
		go tools.MonitorCase(ticker, wetres)
		//defer close(tools.ToolsChan)

		cbox := container.NewHBox(
			widget.NewButton("base64加密", func() {
				tools.Base64Origin = wetin.Text
				encoded := base64.StdEncoding.EncodeToString([]byte(tools.Base64Origin))
				log.Println("base64-encoded:", encoded)
				wetout.SetText(encoded)
				wetres.SetText(encoded)
			}),
			widget.NewButton("base64解密", func() {
				tools.Base64Encode = wetout.Text
				log.Println("base64-decoded-input:", wetout.Text)
				decoded, _ := base64.StdEncoding.DecodeString(tools.Base64Encode)
				log.Println("base64-decoded:", string(decoded))
				wetin.SetText(string(decoded))
				wetres.SetText(string(decoded))
			}),
			widget.NewButton("sha256加密", func() {
				log.Println("sha256-encoded-input:", wetin.Text)
				h := sha256.Sum256([]byte(wetin.Text))
				log.Println("base64-decoded:", fmt.Sprintf("%x", h))
				wetout.SetText(fmt.Sprintf("%x", h))
				wetres.SetText(fmt.Sprintf("%x", h))
			}),
			widget.NewButton("md5", func() {
				log.Println("md5-input:", wetin.Text)
				sum := md5.Sum([]byte(wetin.Text))
				log.Println("md5-output:", fmt.Sprintf("%x", sum))
				wetout.SetText(fmt.Sprintf("%x", sum))
				wetres.SetText(fmt.Sprintf("%x", sum))
			}),
			widget.NewButton("文件md5", func() {
				dia := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
					if err != nil {
						log.Println("newOpenFileFolder:", err.Error())
						return
					}
					if uc == nil || uc.URI() == nil {
						return
					}
					filePath := uc.URI().Path()
					log.Println("path:", filePath)
					fileMd5, md5Err := tools.CalculateMD5(filePath)
					if md5Err != nil {
						wetout.SetText("生成文件md5失败啦, err:" + md5Err.Error())
						return
					} else {
						wetout.SetText(uc.URI().Name() + " - 文件md5:" + fileMd5)
						wetres.SetText(fileMd5)
					}
					log.Println("path:", fileMd5)
				}, tw)
				dia.Show()
			}),
			widget.NewButton("urlEncode", func() {
				log.Println("urlEncode-input:", wetin.Text)
				sum := url.QueryEscape(wetin.Text)
				log.Println("urlEncode-output:,", sum)
				wetout.SetText(sum)
				wetres.SetText(sum)
			}),
			widget.NewButton("urlDecode", func() {
				log.Println("urlDecode-input:", wetout.Text)
				out, _ := url.QueryUnescape(wetout.Text)
				log.Println("urlDecode-output:", out)
				wetout.SetText(out)
				wetres.SetText(out)
			}),
			widget.NewButton("时间戳", func() {
				isTime := tools.IsTimeFormat(wetin.Text, "2006-01-02 15:04:05")
				if !isTime {
					wetin.SetText("")
				}
				wetin.PlaceHolder = "传入时间格式yyyy-MM-dd HH:mm:ss可输出某时刻的时间戳，未传入这输出当前时间戳!!"
				log.Println("时间戳-input:", wetin.Text)
				timestamp, _ := tools.GetTimestamp(wetin.Text, "2006-01-02 15:04:05")
				log.Println("时间戳-output:", strconv.FormatInt(timestamp, 10))
				if isTime {
					wetout.SetText("时间戳： " + strconv.FormatInt(timestamp, 10))
					wetres.SetText("时间戳： " + strconv.FormatInt(timestamp, 10))
				} else {
					wetout.SetText("当前时间戳： " + strconv.FormatInt(timestamp, 10))
					wetres.SetText("当前时间戳： " + strconv.FormatInt(timestamp, 10))
				}
			}),
		)
		cboxImg := container.NewHBox(
			widget.NewButton("二维码解码", func() {
				dia := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
					if err != nil {
						log.Println("newOpenFileFolder:", err.Error())
						return
					}
					if uc == nil || uc.URI() == nil {
						return
					}
					suffix := uc.URI().Extension()
					log.Println("suffix:", suffix)
					if !tools.IsImg(suffix) {
						wetin.SetText("请打开图片格式文件哦！")
						log.Println("notImg-suffix:", suffix)
						return
					}
					path := uc.URI().String()
					wetin.SetText("打开二维码路径：" + path)
					log.Println("path:", path)
					path = strings.Replace(path, "file://", "", 1)
					res := tools.ReadQRCode(path)
					log.Println("res:", res)
					wetout.SetText("解析 成功-- 解析内容为： " + res)
					wetres.SetText(res)
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
							qrwin.Close()
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
				log.Println("sha512-input:", wetin.Text)
				hashedBytes := sha512.Sum512([]byte(wetin.Text))
				hashedString := fmt.Sprintf("%x", hashedBytes)
				log.Println("sha512-output:", hashedString)
				wetout.SetText(hashedString)
				wetres.SetText(hashedString)
			}),
			widget.NewButton("随机密码", func() {
				wetin.PlaceHolder = "可以输入要生成的随机密码长度，不输入默认取8"
				pwd := ""
				if !tools.IsNumeric(wetin.Text) {
					pwd, _ = tools.GenerateRandomPassword(8)
				} else {
					leng, _ := strconv.Atoi(wetin.Text)
					if leng > 1700 {
						wetout.SetText("/**\n大哥，过分了，\n这么长的密码你能记得住吗，\n收手吧阿祖\n**/")
						return
					}
					pwd, _ = tools.GenerateRandomPassword(leng)
				}
				wetout.SetText("生成随记密码:" + pwd)
				wetres.SetText(pwd)
				log.Println("password-output:", pwd)
			}),
			widget.NewButton("json美化", func() {
				// 解析JSON数据到一个map或结构体
				var parsedData map[string]interface{}
				if err := json.Unmarshal([]byte(wetin.Text), &parsedData); err != nil {
					log.Println("JSON解析失败:", err.Error())
					wetout.SetText("json解析失败：" + err.Error())
					return
				}
				prettyJson, err := tools.PrettyPrintJSON(parsedData)
				if err != nil {
					wetout.SetText("json解析失败")
					return
				}
				log.Println("json美化：" + prettyJson)
				wetout.SetText(prettyJson)
				wetres.SetText(prettyJson)
			}),
			widget.NewButton("字数统计", func() {
				log.Println("字数统计-input:", wetin.Text)
				num := tools.CountValidWords(wetin.Text)
				log.Println("字数统计-output:%d", num)
				wetout.SetText("字数统计(不统计空格):" + strconv.Itoa(num))
				wetres.SetText(strconv.Itoa(num))
			}),
			widget.NewCheck("转大写", func(value bool) {
				log.Println("大小写：flag:", value)
				tools.IsLowerCase = value // 设置标志来表示是否要监听剪贴板
			}),
		)
		cboxPro := container.NewHBox(
			widget.NewButton("二维码解析增强版", func() {
				// 文件不存在，提示可以下载
				toolName := settings.Conf.ToolName
				toolUrl := settings.Conf.DownLoadUrl
				toolMd5 := settings.Conf.ToolMd5
				if !tools.IsFileNameExists(toolName) {
					sew := myApp.NewWindow("下载增强工具")
					wn := widget.NewLabel("是否下载增强工具——大小为465kb")
					toolbtn := widget.NewButton("点击下载工具", func() {
						msg := ""
						if !tools.IsFileNameExists(toolName) {
							if tools.DownloadFile(toolUrl, toolName) {
								msg = "下载工具成功！，关闭窗口即可。\n校验md5:" + toolMd5
							} else {
								msg = "下载工具失败，请检查网络后重试, \n或者打开项目git地址检查是否是最新版本 \n项目地址：https://github.com/xs-alpha/xs_camel！"
							}
						} else {
							msg = "当前工具已经存在啦，咋地，还不放心嘛？\n不放心就校验一下md5吧：" + toolMd5
						}
						dw := myApp.NewWindow("下载反馈！")
						dww := widget.NewLabel(msg)
						dw.Resize(fyne.NewSize(200, 200))
						dw.SetContent(dww)
						dw.Show()

					})
					sew.SetContent(container.NewVBox(wn, toolbtn))
					sew.Show()
				} else {
					// 文件存在
					currentDir, _ := os.Getwd()
					exePath := currentDir + "\\" + toolName
					cmd := exec.Command(exePath)
					cmd.Start()
				}

			}),
			widget.NewButton("随机选择", func() {
				if !tools.NoMoreChooseInfo {
					cow := myApp.NewWindow("提示")
					conl := widget.NewLabel("选择困难症福音!!!")
					com := widget.NewMultiLineEntry()
					com.SetMinRowsVisible(10)
					com.SetText("点击确定后\n在原文位置输入内容，\n多个选择请换行分割")
					cowck := widget.NewCheck("不再提示", func(b bool) {
						tools.NoMoreChooseInfo = b
					})
					wbtn := widget.NewButton("确定", func() {
						cow.Close()
					})
					cow.SetContent(container.NewVBox(conl, com, cowck, wbtn))
					cow.Resize(fyne.NewSize(200, 200))
					cow.Show()
				}

				lines := strings.Split(wetin.Text, "\n")
				log.Println("parse choices:", lines)
				randomString, _ := tools.GetRandomString(lines)
				wetout.SetText("帮您做出了选择：" + randomString)
				wetres.SetText(randomString)
			}),
			widget.NewButton("翻译", func() {
				log.Println("翻译-input:", wetin.Text)
				wetout.SetText("")
				wetres.SetText("")
				if wetin.Text==""{
					wetout.SetText("嘿girl,输入内容呀！")
					return
				}
				tools.TransLate(wetin.Text, wetout,wetres)
			}),
			widget.NewButton("翻译增强", func() {
			    if !tools.FileExists(settings.Conf.ConfigName){
	                bdw:=myApp.NewWindow("请输入百度翻译sdk")
                    bdwappid:=widget.NewMultiLineEntry()
					bdwappid.SetPlaceHolder("请填写百度AppId")
                    bdwsecret:=widget.NewMultiLineEntry()
					bdwsecret.SetPlaceHolder("请填写百度secret")
                    bdwappid.SetMinRowsVisible(2)
                    bdwsecret.SetMinRowsVisible(2)
					wl:=widget.NewLabel("")
                    bdb:=widget.NewButton("保存",func(){
                        if bdwappid.Text!="" && bdwsecret.Text!=""{
                            terr:=tools.WriteConfig(bdwappid.Text,bdwsecret.Text)
                            if terr!=nil{
                                wl.SetText("保存失败，请稍候重试，或者删除bdsdk.ini重试")
                            }else{
                                wetout.SetText("保存成功！！！")
								tools.BaiDuConfig.AppId = bdwappid.Text
								tools.BaiDuConfig.AppSec = bdwsecret.Text
                                bdw.Close()
                            }
                        }else{
							wl.SetText("请填写百度sdk密钥和appId")
                        }
                    })
                    bdw.SetContent(container.NewVBox(bdwappid,bdwsecret,wl,bdb))
                    bdw.Resize(fyne.NewSize(200,200))
                    bdw.Show()
			    }else{
					// 方便的话直接从内存取
                    if tools.BaiDuConfig.AppId=="" || tools.BaiDuConfig.AppSec==""{
						log.Println("内存未找到，读取配置文件")
						// 读取文件
						appID, secret,cerr:=tools.ReadConfig()
						if cerr!=nil{
							terr:=tools.DeleteFile(settings.Conf.ConfigName)
							if terr!=nil{
                                wetout.SetText("读取配置失败，请手动删除，重新配置！！！")
							}
						}
						tools.BaiDuConfig.AppId = appID
						tools.BaiDuConfig.AppSec = secret
					}else{
						log.Println("内存中存在，从内存取配置文件")
					}
					if tools.BaiDuConfig.AppId==""||tools.BaiDuConfig.AppSec==""{
						wetout.SetText("配置错误\n")
						return
					}
					if wetin.Text==""{
						wetout.SetText("嘿gay,输入内容呀\n")
						return
					}
					wg.Add(1)
					go tools.DoTransLate(wetin.Text,wetout,wetres)
					wg.Done()
					

			    }

            }),
		)
		cn := container.NewVBox(wetin, cbox, cboxImg, cboxPro, wetout, wetres)
		tw.SetContent(cn)
		tw.Resize(fyne.NewSize(300, 300))
		tw.Show()
		// 回调函数
		tw.SetOnClosed(func() {
			log.Println("关闭小工具的大小写转换通道")
			close(tools.ToolsChan)
		})

	})
	return toolBtn
}
