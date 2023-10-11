package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"unicode"

	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

var BaiDuConfig = new(BaiduFanyi)
type BaiduFanyi struct {
    AppId, AppSec string
}

func (this *BaiduFanyi) Translate(query, from, to string) (string, error) {
    salt := "1435660288"
    data := url.Values{}
    data.Set("q", query)
    data.Set("salt", salt)
    data.Set("appid", this.AppId)
    data.Set("from", from)
    data.Set("to", to)
    data.Set("sign", this.BuildSign(query, salt))
    res, err := PostForm("http://api.fanyi.baidu.com/api/trans/vip/translate", data)
    fmt.Println(res, err)
    res = gjson.Get(res, "trans_result.0.dst").String()
    res = UrlDecode(res)
    return res, err
}
func (this *BaiduFanyi) BuildSign(query, salt string) string {
    str := fmt.Sprintf("%s%s%s%s", this.AppId, query, salt, this.AppSec)
    ret := Md5(str)
    return ret
}

//发送http post请求数据为form
func PostForm(url string, data url.Values) (string, error) {
    resp, err := http.PostForm(url, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(content), nil
}

//md5加密
func Md5(src string) string {
    m := md5.New()
    m.Write([]byte(src))
    res := hex.EncodeToString(m.Sum(nil))
    return res
}
func UrlDecode(str string) string {
    res, err := url.QueryUnescape(str)
    if err != nil {
        return ""
    }
    return res
}
func isEnglish(text string) bool {
    for _, r := range text {
        if !unicode.Is(unicode.Latin, r) {
            return false
        }
    }
    return true
}
func translateUseSdk(query string)(string,error){
    // baidu := &BaiduFanyi{
    //     AppId:  "2022070xxxxxx5",
    //     AppSec: "xxxxxxx",
    // }
    baidu:=BaiDuConfig
    if isEnglish(query){
         return baidu.Translate(query, "en", "zh")
    }else{
        return baidu.Translate(query, "zh", "en")
    }
}

func DoTransLate(query string,wetout, wetres *widget.Entry ){
    res,transErr:=translateUseSdk(query)
    if transErr!=nil{
        wetout.SetText("翻译失败，请检查网络、appid、secret")
    }else{
        wetout.SetText("【百度翻译Pro:】:\n"+res)
        wetres.SetText("【百度翻译Pro:】:\n"+res)
    }
}

